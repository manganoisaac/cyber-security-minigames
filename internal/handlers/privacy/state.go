package privacy

import (
	"cyber-security-minigames/data"
	"cyber-security-minigames/internal/session"
	"cyber-security-minigames/types"
	"net/http"
	"sync"
)

const stateKey = "privacy"

type IPrivacyGameState interface {
	session.IState
	Posts() []types.Post
	ClickPost(id int) types.Post
	AllPostsAnalysed() bool
}

type gameState struct {
	mu    sync.Mutex
	posts []types.Post
}

func newGameState() IPrivacyGameState {
	return &gameState{
		posts: data.GetPostData(),
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

	posts := make([]types.Post, len(s.posts))
	copy(posts, s.posts)
	return posts
}

func (s *gameState) ClickPost(id int) types.Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.posts {
		if s.posts[i].Id == id {
			s.posts[i].Analysed = true
			if s.posts[i].Clue != nil {
				s.posts[i].Clue.AddedToFolder = true
			}
			return s.posts[i]
		}
	}
	return types.Post{}
}

func (s *gameState) AllPostsAnalysed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.posts) == 0 {
		return false
	}

	for _, post := range s.posts {
		if !post.Analysed {
			return false
		}
	}

	return true
}
