package email_client

import (
	"fmt"
	"sync"

	"cyber-security-minigames/data"
	"cyber-security-minigames/types"
	"cyber-security-minigames/utils/emailtext"
	"cyber-security-minigames/views/components/email"
	"cyber-security-minigames/views/components/ui"
	"cyber-security-minigames/views/pages"

	"github.com/a-h/templ"
)

// submitResult is the outcome of grading a phrase submission: the feedback
// message and, on success, the page the player is sent to.
type submitResult struct {
	message string
	nextURL string
	correct bool
}

// checkResult is the outcome of checking an email: the verdict and, when it
// was malware, the page that starts level 2.
type checkResult struct {
	verdict string
	url     string
}

// IEmailService owns the email minigame's orchestration. Methods receive only
// the parameters they need; the HTTP handlers decode requests and render the
// returned components.
type IEmailService interface {
	Page(st IEmailGameState) templ.Component
	Instructions(st IEmailGameState) templ.Component
	Select(st IEmailGameState, em types.Email) []templ.Component
	TogglePhrase(st IEmailGameState, em types.Email, phraseIdx int, sel bool) ([]templ.Component, error)
	Submit(st IEmailGameState, em types.Email) submitResult
	RevealSelected(st IEmailGameState, em types.Email) []templ.Component
	Check(st IEmailGameState, em types.Email) checkResult
}

// Service implements IEmailService, delegating submission grading to an
// injected GradingService.
type Service struct {
	grading GradingService
}

// newService is the private constructor for Service.
func newService(grading GradingService) *Service {
	return &Service{grading: grading}
}

var (
	serviceOnce   sync.Once
	serviceGlobal IEmailService
)

// EmailService returns the app-wide email service singleton, constructed once
// on first use.
func EmailService() IEmailService {
	serviceOnce.Do(func() { serviceGlobal = newService(NewGradingService()) })
	return serviceGlobal
}

func (s *Service) Page(st IEmailGameState) templ.Component {
	active, _, ok := data.EmailByID(st.ActiveEmail())
	if !ok || st.Level() >= 2 {
		active, _, _ = data.PhishingEmail()
		st.SetActiveEmail(active.ID)
	}
	intro := ""
	if st.Level() == 1 {
		intro = data.InstructionsIntroLevel1
	}
	return pages.EmailClientPage(data.Emails, active, st.SelectedPhrasesMap(active.ID), st.Level(), intro, data.InstructionsForLevel(st.Level()), st.Revealed())
}

// Instructions renders just the instructions modal shown after the player
// dismisses the pre-instruction overview. It is swapped into #modal-root.
func (s *Service) Instructions(st IEmailGameState) templ.Component {
	return ui.Modal(data.InstructionsForLevel(st.Level()), "", false)
}

func (s *Service) Select(st IEmailGameState, em types.Email) []templ.Component {
	old, _, oldOK := data.EmailByID(st.ActiveEmail())
	st.SetActiveEmail(em.ID)

	parts := []templ.Component{email.EmailDetails(em, st.SelectedPhrasesMap(em.ID), st.Level(), st.Revealed())}
	if oldOK && old.ID != em.ID {
		parts = append(parts, email.EmailItem(old, false, true))
	}
	return append(parts, email.EmailItem(em, true, true))
}

func (s *Service) TogglePhrase(st IEmailGameState, em types.Email, phraseIdx int, sel bool) ([]templ.Component, error) {
	phrase, ok := emailtext.PhraseAt(em, phraseIdx)
	if !ok {
		return nil, fmt.Errorf("invalid phrase index %d", phraseIdx)
	}
	if sel {
		st.SelectPhrase(em.ID, phraseIdx)
	} else {
		st.DeselectPhrase(em.ID, phraseIdx)
	}

	correctSet := emailtext.CorrectSet(em)
	phraseCmp := email.EmailPhrase(phrase, sel, phraseIdx, em.ID, correctSet[phraseIdx], st.Revealed(), false)
	submitCmp := email.EmailSubmit("Submit", st.SelectedCount(em.ID), true, em.ID)
	return []templ.Component{phraseCmp, submitCmp}, nil
}

func (s *Service) Submit(st IEmailGameState, em types.Email) submitResult {
	res := s.grading.Grade(st, em)
	st.SetRevealed(true)
	return res
}

// RevealSelected re-renders every currently-selected phrase with its
// correctness now showing, as out-of-band fragments the Submit response adds
// alongside its result modal - selection itself never gives this away, only
// pressing Submit does.
func (s *Service) RevealSelected(st IEmailGameState, em types.Email) []templ.Component {
	correctSet := emailtext.CorrectSet(em)
	var parts []templ.Component
	for _, idx := range st.SelectedPhraseIndexes(em.ID) {
		phrase, ok := emailtext.PhraseAt(em, idx)
		if !ok {
			continue
		}
		parts = append(parts, email.EmailPhrase(phrase, true, idx, em.ID, correctSet[idx], true, true))
	}
	return parts
}

func (s *Service) Check(st IEmailGameState, em types.Email) checkResult {
	if em.IsMalware {
		st.SetLevel(2)
		st.SetActiveEmail(em.ID)
		return checkResult{
			verdict: fmt.Sprintf("This email is suspicious! %s", em.Reason),
			url:     "/email-client",
		}
	}
	return checkResult{
		verdict: fmt.Sprintf("This email is not suspicious. %s", em.Reason),
	}
}
