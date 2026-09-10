package password

import (
	"net/http"
	"slices"
	"sync"

	"cyber-security-minigames/data"
	"cyber-security-minigames/internal/session"
	"cyber-security-minigames/types"
)

// stateKey is the name of this minigame's state slot inside the shared session.
const stateKey = "password"

// IPasswordGameState is the state contract for the password minigame. Only
// Reset is shared with other games; every other operation is game-specific.
type IPasswordGameState interface {
	session.IState
	Guesses() []types.Guess
	AddGuess(guess string)
	Backspace()
	CorrectPhrases() []string
	View() types.PasswordView
	SetView(view types.PasswordView)
	AddHint(hint string)
	GetHints() []string
	AllHintsCollected() bool
	Level() int
}

// gameState is the per-user state for the password minigame. The server owns
// every piece of state; the HTML rendered to the client is just a projection.
type gameState struct {
	level          int
	hints          []string
	numHints       int
	mu             sync.Mutex
	guesses        []types.Guess
	correctPhrases []string
	view           types.PasswordView
}

// newGameState returns the concrete state; callers receive it through
// IPasswordGameState.
func newGameState() *gameState {
	return &gameState{
		correctPhrases: data.PasswordPhrases,
		level:          1,
		// 5 real hint-granting objects (jersey, calendar, dog photo, drawing,
		// ID tag) - the trash can no longer counts, see AddHint's objectId
		// case 3 in service.go.
		numHints: 5,
		guesses:  []types.Guess{},
	}
}

// stateFor returns the password minigame's state for a session, creating it on
// first use. The state is stored in the shared session so it survives
// switching to other minigames and back.
func stateFor(s session.ISession) IPasswordGameState {
	if st, ok := s.State(stateKey); ok {
		return st.(IPasswordGameState)
	}
	gs := newGameState()
	s.SetState(stateKey, gs)
	return gs
}

// stateForRequest returns the password minigame's state for the request's
// session, as resolved into the context by session.Middleware.
func stateForRequest(r *http.Request) IPasswordGameState {
	return stateFor(session.SessionFrom(r))
}

func (s *gameState) AddHint(hint string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !slices.Contains(s.hints, hint) {
		s.hints = append(s.hints, hint)
	}
}

func (s *gameState) IncrementLevel() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.level = s.level + 1
}

func (s *gameState) Level() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.level
}

func (s *gameState) Guesses() []types.Guess {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.guesses
}

func (s *gameState) AddGuess(guess string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.guesses) >= len(s.correctPhrases) {
		return
	}
	var wordStatus types.GuessStatus
	if s.correctPhrases[len(s.guesses)] == guess {
		wordStatus = types.GuessCorrect
	} else if slices.Contains(s.correctPhrases, guess) {
		wordStatus = types.GuessWrongPosition
	} else {
		wordStatus = types.GuessIncorrect
	}
	guessInsert := types.Guess{
		Word:   guess,
		Status: wordStatus,
	}
	s.guesses = append(s.guesses, guessInsert)
}

func (s *gameState) Backspace() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.guesses) > 0 {
		s.guesses = s.guesses[:len(s.guesses)-1]
	}
}

func (s *gameState) CorrectPhrases() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.correctPhrases
}

func (s *gameState) View() types.PasswordView {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.view
}

func (s *gameState) SetView(view types.PasswordView) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.view = view
}

// Reset clears this minigame's per-run state.
func (s *gameState) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.guesses = nil
}

func (s *gameState) GetHints() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hints
}

func (s *gameState) AllHintsCollected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := false
	if s.numHints <= len(s.hints) {
		res = true
		s.level = s.level + 1
	}
	return res
}
