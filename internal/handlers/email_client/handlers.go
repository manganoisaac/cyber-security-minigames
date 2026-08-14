package email_client

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"cyber-security-minigames/data"
	"cyber-security-minigames/utils/emailtext"
	"cyber-security-minigames/views/components/email"
	"cyber-security-minigames/views/components/ui"
	"cyber-security-minigames/views/pages"
)

func HandlePage(w http.ResponseWriter, r *http.Request) {
	state := sessions.get(w, r)
	active, _, ok := data.EmailByID(state.activeEmail())
	if !ok || state.level() >= 2 {
		active, _, _ = data.PhishingEmail()
		state.setActiveEmail(active.ID)
	}
	pages.EmailClientPage(data.Emails, active, state.selectedPhrasesMap(active.ID), state.level(), data.InstructionsForLevel(state.level())).Render(r.Context(), w)
}

func HandleSelect(w http.ResponseWriter, r *http.Request) {
	em, _, ok := data.EmailByID(r.PathValue("id"))
	if !ok {
		http.Error(w, "invalid email id", http.StatusBadRequest)
		return
	}
	state := sessions.get(w, r)
	old, _, oldOK := data.EmailByID(state.activeEmail())
	state.setActiveEmail(em.ID)

	email.EmailDetails(em, state.selectedPhrasesMap(em.ID), state.level()).Render(r.Context(), w)
	if oldOK && old.ID != em.ID {
		email.EmailItem(old, false, true).Render(r.Context(), w)
	}
	email.EmailItem(em, true, true).Render(r.Context(), w)
}

func HandlePhrase(w http.ResponseWriter, r *http.Request) {
	em, _, ok := data.EmailByID(r.PathValue("emailID"))
	if !ok {
		http.Error(w, "invalid email id", http.StatusBadRequest)
		return
	}
	phraseIndex, err := strconv.Atoi(r.PathValue("phraseIndex"))
	if err != nil {
		http.Error(w, "phraseIndex is not a number", http.StatusBadRequest)
		return
	}
	sel := r.PathValue("select") == "select"

	state := sessions.get(w, r)
	if sel {
		state.selectPhrase(em.ID, phraseIndex)
	} else {
		state.deselectPhrase(em.ID, phraseIndex)
	}

	phrase, ok := emailtext.PhraseAt(em, phraseIndex)
	if !ok {
		http.Error(w, "invalid phrase index", http.StatusBadRequest)
		return
	}
	correctSet := emailtext.CorrectSet(em)
	email.EmailPhrase(phrase, sel, phraseIndex, em.ID, correctSet[phraseIndex]).Render(r.Context(), w)
	email.EmailSubmit("Submit", state.selectedCount(em.ID), true, em.ID).Render(r.Context(), w)
}

func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	em, _, ok := data.EmailByID(r.PathValue("emailID"))
	if !ok {
		http.Error(w, "invalid email id", http.StatusBadRequest)
		return
	}

	state := sessions.get(w, r)
	selected := state.selectedPhraseIndexes(em.ID)

	correct := make(map[int]bool, len(em.CorrectPhraseIndexes))
	for _, idx := range em.CorrectPhraseIndexes {
		correct[idx] = true
	}

	missing, wrong := 0, 0
	for _, idx := range em.CorrectPhraseIndexes {
		if !selectedContains(selected, idx) {
			missing++
		}
	}
	for _, idx := range selected {
		if !correct[idx] {
			wrong++
		}
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

	url := ""
	if missing == 0 && wrong == 0 {
		url = "/chapter-select"
		state.reset()
	}
	ui.Modal(message, url).Render(r.Context(), w)
}

func selectedContains(selected []int, idx int) bool {
	for _, s := range selected {
		if s == idx {
			return true
		}
	}
	return false
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

func HandleCheck(w http.ResponseWriter, r *http.Request) {
	em, _, ok := data.EmailByID(r.PathValue("emailID"))
	if !ok {
		http.Error(w, "invalid email id", http.StatusBadRequest)
		return
	}
	state := sessions.get(w, r)
	verdict := fmt.Sprintf("This email is not suspicious. %s", em.Reason)
	var url string
	if em.IsMalware {
		state.setLevel(2)
		state.setActiveEmail(em.ID)
		verdict = fmt.Sprintf("This email is suspicious! %s", em.Reason)
		url = "/email-client"
	}
	ui.Modal(verdict, url).Render(r.Context(), w)
}

func HandleModalClose(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
