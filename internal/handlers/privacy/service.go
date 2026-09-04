package privacy

import (
	"cyber-security-minigames/types"
	"cyber-security-minigames/utils/privacyutils"
	"cyber-security-minigames/views/components/privacy"
	"cyber-security-minigames/views/components/ui"
	"cyber-security-minigames/views/pages"
	"sync"

	"github.com/a-h/templ"
)

type IPrivacyService interface {
	Page(IPrivacyGameState) templ.Component
	ClickPost(id int, st IPrivacyGameState) []templ.Component
}

type Service struct{}

func newService() IPrivacyService {
	return &Service{}
}

var (
	serviceOnce   sync.Once
	serviceGlobal IPrivacyService
)

func PrivacyService() IPrivacyService {
	serviceOnce.Do(func() { serviceGlobal = newService() })
	return serviceGlobal
}

func (s *Service) Page(st IPrivacyGameState) templ.Component {
	data := st.Posts()
	postData := privacyutils.PostDataToUI(data)
	return pages.PrivacyPage(postData.Posts, postData.Clues)
}

func (s *Service) ClickPost(id int, st IPrivacyGameState) []templ.Component {
	post := st.ClickPost(id)
	postUI := types.PostUI{
		Id:      post.Id,
		Author:  post.Author,
		Content: post.Content,
	}
	if post.Clue != nil {
		postUI.AddedToFolder = post.Clue.AddedToFolder
	}

	posts := st.Posts()
	privacyUI := privacyutils.PostDataToUI(posts)
	components := []templ.Component{
		privacy.PostItem(postUI),
		privacy.ClueList(privacyUI.Clues, true),
	}
	if st.AllCluesFound(len(privacyUI.Clues)) {
		components = append(components, ui.Modal(
			"congratulations you found all of the clues",
			"/chapter-select",
			true,
		))
	}
	return components
}
