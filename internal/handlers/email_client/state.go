package email_client

import (
	"net/http"
	"sync"

	"cyber-security-minigames/internal/session"
)

// stateKey is the name of this minigame's state slot inside the shared session.
const stateKey = "email-client"

// IEmailGameState is the state contract for the email minigame. Only Reset
// is shared with other games; every other operation is email-specific.
type IEmailGameState interface {
	session.IState
	ActiveEmail() string
	SetActiveEmail(id string)
	SelectPhrase(emailID string, phraseIdx int)
	DeselectPhrase(emailID string, phraseIdx int)
	SelectedCount(emailID string) int
	SelectedPhraseIndexes(emailID string) []int
	SelectedPhrasesMap(emailID string) map[int]bool
	HasSelected(emailID string, phraseIdx int) bool
	Level() int
	SetLevel(level int)
	Revealed() bool
	SetRevealed(revealed bool)
}

// gameState is the per-user state for the email minigame. The server owns
// every piece of state; the HTML rendered to the client is just a projection.
type gameState struct {
	mu              sync.Mutex
	activeID        string
	selectedPhrases map[string]map[int]bool // email id -> (phrase index -> selected)
	gameLevel       int                     // 1 = find the phishing email, 2 = phrase hunt
	revealed        bool                    // whether the last Submit's per-phrase correctness is showing
}

// newGameState returns the concrete state; callers receive it through
// IEmailGameState.
func newGameState() *gameState {
	return &gameState{
		selectedPhrases: make(map[string]map[int]bool),
		gameLevel:       1,
		activeID:        "mobile-plan-renewal-notice",
	}
}

// stateFor returns the email minigame's state for a session, creating it on
// first use. The state is stored in the shared session so it survives
// switching to other minigames and back.
func stateFor(s session.ISession) IEmailGameState {
	if st, ok := s.State(stateKey); ok {
		return st.(IEmailGameState)
	}
	gs := newGameState()
	s.SetState(stateKey, gs)
	return gs
}

// stateForRequest returns the email minigame's state for the request's
// session, as resolved into the context by session.Middleware.
func stateForRequest(r *http.Request) IEmailGameState {
	return stateFor(session.SessionFrom(r))
}

func (s *gameState) ActiveEmail() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.activeID
}

func (s *gameState) SetActiveEmail(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeID = id
}

func (s *gameState) SelectPhrase(emailID string, phraseIdx int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	phrases, ok := s.selectedPhrases[emailID]
	if !ok {
		phrases = make(map[int]bool)
		s.selectedPhrases[emailID] = phrases
	}
	phrases[phraseIdx] = true
	// Changing the selection re-hides last submission's reveal - correctness
	// only shows again once the player submits this new selection.
	s.revealed = false
}

func (s *gameState) DeselectPhrase(emailID string, phraseIdx int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if phrases, ok := s.selectedPhrases[emailID]; ok {
		delete(phrases, phraseIdx)
	}
	s.revealed = false
}

func (s *gameState) SelectedCount(emailID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.selectedPhrases[emailID])
}

func (s *gameState) SelectedPhraseIndexes(emailID string) []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	phrases := s.selectedPhrases[emailID]
	indexes := make([]int, 0, len(phrases))
	for idx := range phrases {
		indexes = append(indexes, idx)
	}
	return indexes
}

func (s *gameState) HasSelected(emailID string, phraseIdx int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	phrases, ok := s.selectedPhrases[emailID]
	if !ok {
		return false
	}
	return phrases[phraseIdx]
}

func (s *gameState) Level() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.gameLevel
}

func (s *gameState) SetLevel(level int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gameLevel = level
}

// Revealed reports whether the last Submit's per-phrase correctness should
// currently be shown (cleared again the moment the player changes their
// selection, so it only ever reflects the most recent Submit).
func (s *gameState) Revealed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.revealed
}

func (s *gameState) SetRevealed(revealed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.revealed = revealed
}

// Reset clears this minigame's per-run state: selected phrases, active email,
// level, and reveal. Other minigames' session state is left untouched.
func (s *gameState) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.selectedPhrases = make(map[string]map[int]bool)
	s.activeID = ""
	s.gameLevel = 1
	s.revealed = false
}

func (s *gameState) SelectedPhrasesMap(emailID string) map[int]bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	phrases := s.selectedPhrases[emailID]
	out := make(map[int]bool, len(phrases))
	for idx, sel := range phrases {
		out[idx] = sel
	}
	return out
}
