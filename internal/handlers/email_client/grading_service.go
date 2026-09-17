package email_client

import (
	"fmt"
	"log"

	"cyber-security-minigames/types"
)

// GradingService computes submission feedback for the email phrase hunt.
type GradingService interface {
	Grade(st IEmailGameState, em types.Email) submitResult
}

type gradingService struct{}

// NewGradingService returns the default grading implementation.
func NewGradingService() GradingService {
	return gradingService{}
}

func (gradingService) Grade(st IEmailGameState, em types.Email) submitResult {
	selected := st.SelectedPhraseIndexes(em.ID)

	correct := make(map[int]bool, len(em.CorrectPhraseIndexes))
	for _, idx := range em.CorrectPhraseIndexes {
		correct[idx] = true
	}

	missing, wrong := 0, 0
	for _, idx := range em.CorrectPhraseIndexes {
		if !st.HasSelected(em.ID, idx) {
			missing++
		}
	}
	for _, idx := range selected {
		if !correct[idx] {
			wrong++
		}
	}

	score := 0
	for _, idx := range selected {
		if correct[idx] {
			score++
		}
	}
	if missing == 0 && wrong == 0 {
		log.Printf("email-submit: correct (%d/%d) for %q", score, len(em.CorrectPhraseIndexes), em.Title)
	} else {
		log.Printf("email-submit: incorrect (%d/%d) for %q", score, len(em.CorrectPhraseIndexes), em.Title)
	}

	var message string
	switch {
	case missing == 0 && wrong == 0:
		message = "Good job! You found all the suspicious phrases."
	case missing > 0 && wrong == 0:
		message = fmt.Sprintf("There are still %d more suspicious %s you haven't selected.", missing, pluralPhrase(missing))
	case missing == 0 && wrong > 0:
		message = fmt.Sprintf("You selected %d %s that %s not suspicious.", wrong, pluralPhrase(wrong), wereAux(wrong))
	default:
		message = fmt.Sprintf("There are still %d suspicious %s you haven't selected, and you selected %d %s that %s not suspicious.", missing, pluralPhrase(missing), wrong, pluralPhrase(wrong), wereAux(wrong))
	}

	nextURL := ""
	if missing == 0 && wrong == 0 {
		nextURL = "/whwl/phishing"
	}

	return submitResult{
		message: message,
		nextURL: nextURL,
		correct: missing == 0 && wrong == 0,
	}
}

func pluralPhrase(n int) string {
	if n == 1 {
		return "phrase"
	}
	return "phrases"
}

func wereAux(n int) string {
	if n == 1 {
		return "is"
	}
	return "are"
}
