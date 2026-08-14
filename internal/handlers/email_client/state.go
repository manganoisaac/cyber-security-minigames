package email_client

import (
	"net/http"
	"sync"

	"cyber-security-minigames/internal/session"
)

// stateKey is the name of this minigame's state slot inside the shared session.
const stateKey = "email-client"

// gameState is the per-user state for the email minigame. The server owns
// every piece of state; the HTML rendered to the client is just a projection.
type gameState struct {
	mu              sync.Mutex
	activeID        string
	selectedPhrases map[string]map[int]bool // email id -> (phrase index -> selected)
	gameLevel       int                     // 1 = find the phishing email, 2 = phrase hunt
}

func newGameState() *gameState {
	return &gameState{
		selectedPhrases: make(map[string]map[int]bool),
		gameLevel:       1,
	}
}

// stateFor returns the email minigame's state for a session, creating it on
// first use. The state is stored in the shared session so it survives
// switching to other minigames and back.
func stateFor(s *session.Session) *gameState {
	if st, ok := s.State(stateKey); ok {
		return st.(*gameState)
	}
	gs := newGameState()
	s.SetState(stateKey, gs)
	return gs
}

// getState returns the email minigame's state for the request's session.
func getState(w http.ResponseWriter, r *http.Request) *gameState {
	return stateFor(session.Default.Get(w, r))
}

func (s *gameState) activeEmail() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.activeID
}

func (s *gameState) setActiveEmail(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeID = id
}

func (s *gameState) selectPhrase(emailID string, phraseIdx int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	phrases, ok := s.selectedPhrases[emailID]
	if !ok {
		phrases = make(map[int]bool)
		s.selectedPhrases[emailID] = phrases
	}
	phrases[phraseIdx] = true
}

func (s *gameState) deselectPhrase(emailID string, phraseIdx int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if phrases, ok := s.selectedPhrases[emailID]; ok {
		delete(phrases, phraseIdx)
	}
}

func (s *gameState) selectedCount(emailID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.selectedPhrases[emailID])
}

func (s *gameState) selectedPhraseIndexes(emailID string) []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	phrases := s.selectedPhrases[emailID]
	indexes := make([]int, 0, len(phrases))
	for idx := range phrases {
		indexes = append(indexes, idx)
	}
	return indexes
}

func (s *gameState) level() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.gameLevel
}

func (s *gameState) setLevel(level int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gameLevel = level
}

// reset clears this minigame's per-run state: selected phrases, active email,
// and level. Other minigames' session state is left untouched.
func (s *gameState) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.selectedPhrases = make(map[string]map[int]bool)
	s.activeID = ""
	s.gameLevel = 1
}

func (s *gameState) selectedPhrasesMap(emailID string) map[int]bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	phrases := s.selectedPhrases[emailID]
	out := make(map[int]bool, len(phrases))
	for idx, sel := range phrases {
		out[idx] = sel
	}
	return out
}
