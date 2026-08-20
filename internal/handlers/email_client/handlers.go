package email_client

import (
	"net/http"
	"strconv"

	"cyber-security-minigames/data"
	"cyber-security-minigames/types"
	"cyber-security-minigames/views/components/ui"
)

// IEmailHandler is the HTTP handler contract for the email minigame.
type IEmailHandler interface {
	Page(w http.ResponseWriter, r *http.Request)
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
	if res.correct {
		st.Reset()
	}
	ui.Modal(res.message, res.nextURL, false).Render(r.Context(), w)
}

func (h *emailHandler) Check(w http.ResponseWriter, r *http.Request) {
	em, ok := h.emailFromPath(w, r, "emailID")
	if !ok {
		return
	}
	res := h.svc.Check(stateForRequest(r), em)
	ui.Modal(res.verdict, res.url, false).Render(r.Context(), w)
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
