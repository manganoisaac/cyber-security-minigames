package privacy

import (
	"cyber-security-minigames/data"
	"cyber-security-minigames/utils/privacyutils"
	"cyber-security-minigames/views/components/privacy"
	"cyber-security-minigames/views/components/ui"
	"cyber-security-minigames/views/pages"
	"sync"

	"github.com/a-h/templ"
)

type IPrivacyService interface {
	Page(IPrivacyGameState) templ.Component
	Instructions(IPrivacyGameState) templ.Component
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
	posts := st.Posts()
	postData := privacyutils.PostDataToUI(posts)
	return pages.PrivacyPage(data.PrivacyProfile, postData.Posts, postData.Clues, data.PrivacyIntro)
}

// Instructions renders just the instructions modal shown after the player
// dismisses the pre-instruction overview. It is swapped into #modal-root.
func (s *Service) Instructions(st IPrivacyGameState) templ.Component {
	return ui.Modal(data.PrivacyInstructions, "", false)
}

func (s *Service) ClickPost(id int, st IPrivacyGameState) []templ.Component {
	post := st.ClickPost(id)

	posts := st.Posts()
	privacyUI := privacyutils.PostDataToUI(posts)
	components := []templ.Component{
		privacy.PostItem(privacyutils.PostToUI(post)),
		privacy.ClueList(privacyUI.Clues, true),
	}
	if st.AllCluesFound(len(privacyUI.Clues)) {
		components = append(components, ui.Modal(
			"congratulations you found all of the clues",
			"/whwl/privacy",
			true,
		))
	}
	return components
}
