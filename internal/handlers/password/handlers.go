package password

import (
	"net/http"
	"strconv"

	"cyber-security-minigames/types"
)

// IPasswordHandler is the HTTP handler contract for the password minigame.
type IPasswordHandler interface {
	Page(w http.ResponseWriter, r *http.Request)
	ChangeView(w http.ResponseWriter, r *http.Request)
	Backspace(w http.ResponseWriter, r *http.Request)
	SelectItem(w http.ResponseWriter, r *http.Request)
	AddGuess(w http.ResponseWriter, r *http.Request)
}

// passwordHandler implements IPasswordHandler, delegating to the service.
type passwordHandler struct {
	svc IPasswordService
}

// NewHandler wires the service into the HTTP layer.
func NewHandler(svc IPasswordService) IPasswordHandler {
	return &passwordHandler{svc: svc}
}

// Page renders the password minigame page.
func (h *passwordHandler) Page(w http.ResponseWriter, r *http.Request) {
	h.svc.Page(stateForRequest(r)).Render(r.Context(), w)
}

// SelectItem handles GET /password/selectitem, delegating to the service.
func (h *passwordHandler) SelectItem(w http.ResponseWriter, r *http.Request) {
	h.svc.SelectItem(stateForRequest(r)).Render(r.Context(), w)
}

// AddGuess handles POST /password/guess/{guess}, delegating to the service.
func (h *passwordHandler) AddGuess(w http.ResponseWriter, r *http.Request) {
	guess := r.PathValue("guess")
	if guess == "" {
		http.Error(w, "guess is empty", http.StatusBadRequest)
		return
	}
	for _, part := range h.svc.AddGuess(stateForRequest(r), guess) {
		part.Render(r.Context(), w)
	}
}

// Backspace handles GET /password/backspace, delegating to the service.
func (h *passwordHandler) Backspace(w http.ResponseWriter, r *http.Request) {
	h.svc.Backspace(stateForRequest(r)).Render(r.Context(), w)
}

// ChangeView handles GET /password/changeview/{view_id}, delegating to the service.
func (h *passwordHandler) ChangeView(w http.ResponseWriter, r *http.Request) {
	viewID, err := strconv.Atoi(r.PathValue("view_id"))
	if err != nil {
		http.Error(w, "view_id must be an integer", http.StatusBadRequest)
		return
	}
	view, ok := types.PasswordViewFromInt(viewID)
	if !ok {
		http.Error(w, "unknown view", http.StatusBadRequest)
		return
	}
	h.svc.ChangeView(stateForRequest(r), view).Render(r.Context(), w)
}
