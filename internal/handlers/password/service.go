package password

import (
	"sync"

	"cyber-security-minigames/data"
	"cyber-security-minigames/types"
	"cyber-security-minigames/utils/passwordutil"
	"cyber-security-minigames/views/components/password"
	"cyber-security-minigames/views/components/ui"
	"cyber-security-minigames/views/pages"

	"github.com/a-h/templ"
)

// IPasswordService owns the password minigame's orchestration. Methods receive
// only the parameters they need; the HTTP handlers decode requests and render
// the returned components.
type IPasswordService interface {
	Page(st IPasswordGameState) templ.Component
	ObjectCloseUp(st IPasswordGameState, objectId int) templ.Component
	ChangeView(st IPasswordGameState, view types.PasswordView) templ.Component
	Backspace(st IPasswordGameState) templ.Component
	SelectItem(st IPasswordGameState) templ.Component
	AddGuess(st IPasswordGameState, guess string) []templ.Component
	AddHint(st IPasswordGameState, objectId int, itemId int) []templ.Component
	GetHints(st IPasswordGameState) templ.Component
}

// Service implements IPasswordService.
type Service struct{}

// newService is the private constructor for Service.
func newService() *Service {
	return &Service{}
}

var (
	serviceOnce   sync.Once
	serviceGlobal IPasswordService
)

// PasswordService returns the app-wide password service singleton, constructed
// once on first use.
func PasswordService() IPasswordService {
	serviceOnce.Do(func() { serviceGlobal = newService() })
	return serviceGlobal
}

func (s *Service) Page(st IPasswordGameState) templ.Component {
	if st.Level() == 1 {

		return pages.PasswordPage(st.View(), st.Guesses(), len(st.CorrectPhrases()), data.InstructionsForPassword())
	} else {
		return pages.PasswordLevel2Page(
			// data.PasswordPhrases,
			st.GetHints(),
			st.Guesses(),
		)
	}
}

func (s *Service) Popup(content string) templ.Component {
	return ui.Popup(content)
}

func (s *Service) SelectItem(st IPasswordGameState) templ.Component {
	return nil
}

func (s *Service) AddGuess(st IPasswordGameState, guess string) []templ.Component {
	st.AddGuess(guess)
	results := st.Guesses()
	parts := []templ.Component{password.Guess(results, len(st.CorrectPhrases()))}
	if passwordutil.AllCorrect(results) && len(results) == len(st.CorrectPhrases()) {
		st.Reset()
		parts = append(parts, ui.Modal(data.PasswordWinMessage, "/chapter-select", true))
	}
	return parts
}

func (s *Service) Backspace(st IPasswordGameState) templ.Component {
	st.Backspace()
	return password.Guess(st.Guesses(), len(st.CorrectPhrases()))
}

func (s *Service) ChangeView(st IPasswordGameState, view types.PasswordView) templ.Component {
	st.SetView(view)
	return s.Page(st)
}

func (s *Service) ObjectCloseUp(st IPasswordGameState, objectId int) templ.Component {
	return password.ObjectCloseUp(objectId, 1)
}

func (s *Service) AddHint(st IPasswordGameState, objectId int, itemId int) []templ.Component {
	switch objectId {
	case 1:
		st.AddHint("jersey")
	case 2:
		st.AddHint("calendar")
	case 3:
		// The trash can is a narrative easter egg about Eugene's motive, not
		// a real clue - viewing it must never count toward the puzzle.
		return []templ.Component{s.Popup("Just an old note - doesn't look important.")}
	case 4:
		st.AddHint("rufus")
	case 5:
		st.AddHint("drawing")
	case 6:
		st.AddHint("id")
	default:
		st.AddHint("desk")
	}
	parts := []templ.Component{s.Popup("hint added to notebook")}
	if st.AllHintsCollected() {
		modal := ui.Modal("congratulations, all hints have been collected.", "/password", true)
		parts = append(parts, modal)
	}
	return parts
}

func (s *Service) GetHints(st IPasswordGameState) templ.Component {
	hints := st.GetHints()
	return password.ClueList(hints)
}
