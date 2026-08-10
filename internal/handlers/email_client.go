package handlers

import (
	"cyber-security-minigames/data"
	"cyber-security-minigames/views/components/email"
	"cyber-security-minigames/views/components/ui"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HandleEmailSelect(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 || id >= len(data.Emails) {
		http.Error(w, "invalid email id", http.StatusBadRequest)
		return
	}
	checkedCount, err := strconv.Atoi(r.URL.Query().Get("checkedCount"))
	if err != nil || checkedCount < 0 {
		checkedCount = 0
	}
	component := email.EmailList(data.Emails, id, 0, checkedCount)
	component.Render(r.Context(), w)
}

func HandleEmailWord(w http.ResponseWriter, r *http.Request) {
	emailIndex, err := strconv.Atoi(r.PathValue("emailIndex"))
	if err != nil || emailIndex < 0 || emailIndex >= len(data.Emails) {
		http.Error(w, "invalid email id", http.StatusBadRequest)
		return
	}
	rawWordIndex := r.PathValue("wordIndex")
	wordIndex, err := strconv.Atoi(rawWordIndex)
	if err != nil {
		http.Error(w, "wordIndex is not a number", http.StatusBadRequest)
		return
	}
	rawSelection := r.PathValue("select")
	var selection bool
	if rawSelection == "select" {
		selection = true
	} else {
		selection = false
	}

	word := r.PathValue("word")
	count, err := strconv.Atoi(r.PostFormValue("selectedCount"))
	if err != nil || count < 0 {
		http.Error(w, "invalid count", http.StatusBadRequest)
		return
	}
	if selection {
		count++
	} else if count > 0 {
		count--
	}

	if err := email.EmailWord(word, selection, wordIndex, emailIndex).Render(r.Context(), w); err != nil {
		return
	}
	email.EmailSubmit("Submit", count, true, emailIndex).Render(r.Context(), w)
}

func HandleEmailSubmit(w http.ResponseWriter, r *http.Request) {
	emailIndex, err := strconv.Atoi(r.PathValue("emailIndex"))
	if err != nil || emailIndex < 0 || emailIndex >= len(data.Emails) {
		http.Error(w, "invalid email id", http.StatusBadRequest)
		return
	}

	var selected []int
	for _, token := range strings.Split(r.PostFormValue("selectedIndexes"), ",") {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		idx, err := strconv.Atoi(token)
		if err != nil {
			continue
		}
		selected = append(selected, idx)
	}

	em := data.Emails[emailIndex]
	correct := make(map[int]bool, len(em.CorrectWordIndexes))
	for _, idx := range em.CorrectWordIndexes {
		correct[idx] = true
	}
	matched := len(selected) == len(em.CorrectWordIndexes)
	score := 0
	for _, idx := range selected {
		if correct[idx] {
			score++
		} else {
			matched = false
		}
	}
	if matched {
		log.Printf("email-submit: correct (%d/%d) for %q", score, len(em.CorrectWordIndexes), em.Title)
	} else {
		log.Printf("email-submit: incorrect (%d/%d) for %q", score, len(em.CorrectWordIndexes), em.Title)
	}
	w.WriteHeader(http.StatusNoContent)
}

func HandleEmailChecked(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.PostFormValue("checkedCount"))
	if err != nil || count < 0 {
		count = 0
	}
	msg := fmt.Sprintf("email-checked: %d emails checked", count)
	log.Print(msg)
	ui.Modal(msg).Render(r.Context(), w)
}
