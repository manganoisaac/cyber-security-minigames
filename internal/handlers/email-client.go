package handlers

import (
	"cyber-security-minigames/data"
	"cyber-security-minigames/views/components"
	"net/http"
	"strconv"
)

func HandleEmailSelect(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 || id >= len(data.EmailContent) {
		http.Error(w, "invalid email id", http.StatusBadRequest)
		return
	}
	component := components.EmailList(id)
	component.Render(r.Context(), w)
}

func HandleEmailWord(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		http.Error(w, "id is not a number", http.StatusBadRequest)
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
	component := components.EmailWord(word, selection, id)
	component.Render(r.Context(), w)
}
