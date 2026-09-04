package privacy

import (
	"cyber-security-minigames/data"
	"cyber-security-minigames/internal/session"
	"cyber-security-minigames/types"
	"cyber-security-minigames/utils/privacyutils"
	"net/http"
	"sync"
)

const stateKey = "privacy"

type IPrivacyGameState interface {
	session.IState
	Posts() []types.Post
	ClickPost(id int) types.Post
	AllCluesFound(cluesFound int) bool
}

type gameState struct {
	mu            sync.Mutex
	posts         []types.Post
	numberOfClues int
}

func newGameState() IPrivacyGameState {
	return &gameState{
		posts:         data.GetPostData(),
		numberOfClues: privacyutils.NumberOfClues(data.GetPostData()),
	}
}

func stateFor(s session.ISession) IPrivacyGameState {
	if st, ok := s.State(stateKey); ok {
		return st.(IPrivacyGameState)
	}
	gs := newGameState()
	s.SetState(stateKey, gs)
	return gs
}

func stateForRequest(r *http.Request) IPrivacyGameState {
	return stateFor(session.SessionFrom(r))
}

func (s *gameState) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.posts = data.GetPostData()
}

func (s *gameState) Posts() []types.Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.posts
}

func (s *gameState) ClickPost(id int) types.Post {
	for i := range s.posts {
		if s.posts[i].Id == id {
			if s.posts[i].Clue != nil {
				s.posts[i].Clue.AddedToFolder = true
			}
			return s.posts[i]
		}
	}
	return types.Post{}
}

func (s *gameState) AllCluesFound(cluesFound int) bool {
	if s.numberOfClues <= cluesFound {
		s.Reset()
		return true
	}
	return false
}
