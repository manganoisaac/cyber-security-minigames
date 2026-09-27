package email_client

import (
	"net/http"
	"strconv"

	"cyber-security-minigames/data"
	"cyber-security-minigames/types"
	"cyber-security-minigames/views/components/email"

	"github.com/a-h/templ"
)

// IEmailHandler is the HTTP handler contract for the email minigame.
type IEmailHandler interface {
	Page(w http.ResponseWriter, r *http.Request)
	Instructions(w http.ResponseWriter, r *http.Request)
	Select(w http.ResponseWriter, r *http.Request)
	Phrase(w http.ResponseWriter, r *http.Request)
	Submit(w http.ResponseWriter, r *http.Request)
	Check(w http.ResponseWriter, r *http.Request)
	ModalClose(w http.ResponseWriter, r *http.Request)
}

// emailHandler implements IEmailHandler, delegating to the email service.
type emailHandler struct {
	svc IEmailService
}

// NewHandler wires the email service into the HTTP layer.
func NewHandler(svc IEmailService) IEmailHandler {
	return &emailHandler{svc: svc}
}

func (h *emailHandler) Page(w http.ResponseWriter, r *http.Request) {
	h.svc.Page(stateForRequest(r)).Render(r.Context(), w)
}

func (h *emailHandler) Instructions(w http.ResponseWriter, r *http.Request) {
	h.svc.Instructions(stateForRequest(r)).Render(r.Context(), w)
}

func (h *emailHandler) Select(w http.ResponseWriter, r *http.Request) {
	em, ok := h.emailFromPath(w, r, "id")
	if !ok {
		return
	}
	for _, part := range h.svc.Select(stateForRequest(r), em) {
		part.Render(r.Context(), w)
	}
}

func (h *emailHandler) Phrase(w http.ResponseWriter, r *http.Request) {
	em, phraseIndex, sel, ok := h.phrasePathValues(w, r)
	if !ok {
		return
	}
	parts, err := h.svc.TogglePhrase(stateForRequest(r), em, phraseIndex, sel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, part := range parts {
		part.Render(r.Context(), w)
	}
}

func (h *emailHandler) Submit(w http.ResponseWriter, r *http.Request) {
	em, ok := h.emailFromPath(w, r, "emailID")
	if !ok {
		return
	}
	st := stateForRequest(r)
	res := h.svc.Submit(st, em)
	// Reveal which selected phrases were right/wrong alongside the result
	// modal - only while the game isn't about to reset and redirect away.
	var reveal []templ.Component
	if !res.correct {
		reveal = h.svc.RevealSelected(st, em)
	} else {
		st.Reset()
	}
	if res.correct {
		email.EmailModal("Analysis Complete", res.message, "Continue", res.nextURL, true).Render(r.Context(), w)
	} else {
		email.EmailModal("Check Your Selection", res.message, "Try Again", "", false).Render(r.Context(), w)
	}
	for _, part := range reveal {
		part.Render(r.Context(), w)
	}
}

func (h *emailHandler) Check(w http.ResponseWriter, r *http.Request) {
	em, ok := h.emailFromPath(w, r, "emailID")
	if !ok {
		return
	}
	res := h.svc.Check(stateForRequest(r), em)
	if res.url != "" {
		email.EmailModal("Phishing Email Identified", res.verdict, "Continue to Analysis", res.url, true).Render(r.Context(), w)
	} else {
		email.EmailModal("Not Phishing", res.verdict, "Keep Looking", "", false).Render(r.Context(), w)
	}
}

func (h *emailHandler) ModalClose(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// emailFromPath resolves an email from a path parameter, writing a 400 and
// returning ok=false if the id does not match a known email.
func (h *emailHandler) emailFromPath(w http.ResponseWriter, r *http.Request, key string) (types.Email, bool) {
	em, _, ok := data.EmailByID(r.PathValue(key))
	if !ok {
		http.Error(w, "invalid email id", http.StatusBadRequest)
		return types.Email{}, false
	}
	return em, true
}

// phrasePathValues resolves the email, phrase index and select flag for the
// phrase route, writing a 400 and returning ok=false on any bad input.
func (h *emailHandler) phrasePathValues(w http.ResponseWriter, r *http.Request) (types.Email, int, bool, bool) {
	em, ok := h.emailFromPath(w, r, "emailID")
	if !ok {
		return types.Email{}, 0, false, false
	}
	phraseIndex, err := strconv.Atoi(r.PathValue("phraseIndex"))
	if err != nil {
		http.Error(w, "phraseIndex is not a number", http.StatusBadRequest)
		return types.Email{}, 0, false, false
	}
	return em, phraseIndex, r.PathValue("select") == "select", true
}
