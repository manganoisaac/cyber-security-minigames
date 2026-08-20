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

// IState is the contract every minigame state satisfies. Only behavior the
// session layer can act on generically lives here; everything game-specific
// belongs on a per-game interface (e.g. IEmailGameState).
type IState interface {
	Reset()
}

// ISession is the per-browser session. Each minigame stores its own state in
// the states map, keyed by an opaque minigame name (e.g. "email-client"),
// so sessions survive switching between minigames within the TTL.
type ISession interface {
	State(name string) (IState, bool)
	SetState(name string, st IState)
}

// Session is the in-memory implementation of ISession.
type Session struct {
	mu      sync.Mutex
	sid     string
	touched time.Time
	states  map[string]IState
}

// State returns the minigame state stored under name, if any.
func (s *Session) State(name string) (IState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.states[name]
	return st, ok
}

// SetState stores st under name, replacing any previous value.
func (s *Session) SetState(name string, st IState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[name] = st
}

// IStore maps session id -> ISession. Single in-memory store per process.
type IStore interface {
	Get(w http.ResponseWriter, r *http.Request) ISession
}

// store is the in-memory implementation of IStore.
type store struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

// newStore is the private constructor for the shared session store.
func newStore() IStore {
	return &store{sessions: make(map[string]*Session)}
}

// Get returns the session for the request's cookie, creating one (and setting
// the cookie) if absent. Expired sessions are pruned lazily on each call.
func (s *store) Get(w http.ResponseWriter, r *http.Request) ISession {
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

	sid := s.newID()
	sess := &Session{
		sid:     sid,
		touched: now,
		states:  make(map[string]IState),
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

var (
	storeOnce   sync.Once
	storeGlobal IStore
)

// Store returns the app-wide session store singleton, constructed once on
// first use.
func Store() IStore {
	storeOnce.Do(func() { storeGlobal = newStore() })
	return storeGlobal
}

// newID returns a cryptographically random session id.
func (s *store) newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
