package privacy

import (
	"net/http"
	"strconv"
)

type IPrivacyHandler interface {
	Page(w http.ResponseWriter, r *http.Request)
	ClickPost(w http.ResponseWriter, r *http.Request)
}

type privacyHandler struct {
	svc IPrivacyService
}

func NewHandler(svc IPrivacyService) IPrivacyHandler {
	return &privacyHandler{svc: svc}
}

func (h *privacyHandler) Page(w http.ResponseWriter, r *http.Request) {
	h.svc.Page(stateForRequest(r)).Render(r.Context(), w)

}

func (h *privacyHandler) ClickPost(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, _ := strconv.Atoi(idString)

	for _, component := range h.svc.ClickPost(id, stateForRequest(r)) {
		component.Render(r.Context(), w)
	}
}
