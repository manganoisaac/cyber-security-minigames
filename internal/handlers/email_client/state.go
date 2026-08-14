package email_client

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

const (
	sessionCookie = "cyber_session"
	sessionTTL    = 30 * time.Minute
)

// gameState is the per-user state for the email minigame. The server owns
// every piece of state; the HTML rendered to the client is just a projection.
type gameState struct {
	mu              sync.Mutex
	sid             string
	touched         time.Time
	activeID        string
	selectedPhrases map[string]map[int]bool // email id -> (phrase index -> selected)
	gameLevel       int                     // 1 = find the phishing email, 2 = phrase hunt
}

func newGameState(sid string) *gameState {
	return &gameState{
		sid:             sid,
		selectedPhrases: make(map[string]map[int]bool),
		gameLevel:       1,
	}
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

// reset clears all per-run game state: selected phrases, active email, and level.
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

// sessionStore maps session id -> game state. Single store per process; the
// theory is each user gets a cookie and their own isolated game state.
type sessionStore struct {
	mu       sync.Mutex
	sessions map[string]*gameState
}

func newSessionStore() *sessionStore {
	return &sessionStore{sessions: make(map[string]*gameState)}
}

func (s *sessionStore) get(w http.ResponseWriter, r *http.Request) *gameState {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, state := range s.sessions {
		if now.After(state.touched.Add(sessionTTL)) {
			delete(s.sessions, id)
		}
	}

	if cookie, err := r.Cookie(sessionCookie); err == nil {
		if state, ok := s.sessions[cookie.Value]; ok {
			state.touched = now
			return state
		}
	}

	sid := newSessionID()
	state := newGameState(sid)
	state.touched = now
	s.sessions[sid] = state
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return state
}

var sessions = newSessionStore()

func newSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
