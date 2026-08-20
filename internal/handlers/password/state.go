package password

import (
	"net/http"
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
	Guesses() []string
	AddGuess(guess string)
	Backspace()
	CorrectPhrases() []string
	View() types.PasswordView
	SetView(view types.PasswordView)
}

// gameState is the per-user state for the password minigame. The server owns
// every piece of state; the HTML rendered to the client is just a projection.
type gameState struct {
	mu             sync.Mutex
	guesses        []string
	correctPhrases []string
	view           types.PasswordView
}

// newGameState returns the concrete state; callers receive it through
// IPasswordGameState.
func newGameState() *gameState {
	return &gameState{
		correctPhrases: data.PasswordPhrases,
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

func (s *gameState) Guesses() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.guesses...)
}

func (s *gameState) AddGuess(guess string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.guesses) >= len(s.correctPhrases) {
		return
	}
	s.guesses = append(s.guesses, guess)
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
	return append([]string(nil), s.correctPhrases...)
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
