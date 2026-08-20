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
	ChangeView(st IPasswordGameState, view types.PasswordView) templ.Component
	Backspace(st IPasswordGameState) templ.Component
	SelectItem(st IPasswordGameState) templ.Component
	AddGuess(st IPasswordGameState, guess string) []templ.Component
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
	return pages.PasswordPage(st.View(), passwordutil.GuessResults(st.Guesses(), st.CorrectPhrases()), len(st.CorrectPhrases()), data.InstructionsForPassword())
}

func (s *Service) SelectItem(st IPasswordGameState) templ.Component {
	return nil
}

func (s *Service) AddGuess(st IPasswordGameState, guess string) []templ.Component {
	st.AddGuess(guess)
	results := passwordutil.GuessResults(st.Guesses(), st.CorrectPhrases())
	parts := []templ.Component{password.Guess(results, len(st.CorrectPhrases()))}
	if passwordutil.AllCorrect(results) && len(results) == len(st.CorrectPhrases()) {
		st.Reset()
		parts = append(parts, ui.Modal(data.PasswordWinMessage, "/chapter-select", true))
	}
	return parts
}

func (s *Service) Backspace(st IPasswordGameState) templ.Component {
	st.Backspace()
	return password.Guess(passwordutil.GuessResults(st.Guesses(), st.CorrectPhrases()), len(st.CorrectPhrases()))
}

func (s *Service) ChangeView(st IPasswordGameState, view types.PasswordView) templ.Component {
	st.SetView(view)
	return s.Page(st)
}
