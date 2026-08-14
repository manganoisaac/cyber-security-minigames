package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

const (
	cookieName = "cyber_session"
	sessionTTL = 30 * time.Minute
)

// Session is the per-browser session. Each minigame stores its own state in
// the states map, keyed by an opaque minigame name (e.g. "email-client"),
// so sessions survive switching between minigames within the TTL.
type Session struct {
	mu      sync.Mutex
	sid     string
	touched time.Time
	states  map[string]any
}

// State returns the minigame state stored under name, if any.
func (s *Session) State(name string) (any, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.states[name]
	return st, ok
}

// SetState stores st under name, replacing any previous value.
func (s *Session) SetState(name string, st any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[name] = st
}

// Store maps session id -> Session. Single in-memory store per process.
type Store struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

func NewStore() *Store {
	return &Store{sessions: make(map[string]*Session)}
}

// Get returns the Session for the request's cookie, creating one (and setting
// the cookie) if absent. Expired sessions are pruned lazily on each call.
func (s *Store) Get(w http.ResponseWriter, r *http.Request) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, sess := range s.sessions {
		if now.After(sess.touched.Add(sessionTTL)) {
			delete(s.sessions, id)
		}
	}

	if cookie, err := r.Cookie(cookieName); err == nil {
		if sess, ok := s.sessions[cookie.Value]; ok {
			sess.touched = now
			return sess
		}
	}

	sid := NewID()
	sess := &Session{
		sid:     sid,
		touched: now,
		states:  make(map[string]any),
	}
	s.sessions[sid] = sess
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return sess
}

// Default is the app-wide store all minigames share so they observe the same
// session cookie and per-browser state.
var Default = NewStore()

// NewID returns a cryptographically random session id.
func NewID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
