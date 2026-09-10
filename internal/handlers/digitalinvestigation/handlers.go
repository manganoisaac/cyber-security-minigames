package digitalinvestigation

import (
	"net/http"

	"cyber-security-minigames/views/pages"
)

// IDigitalInvestigationHandler is the HTTP handler contract for the
// digital investigation minigame.
type IDigitalInvestigationHandler interface {
	Page(w http.ResponseWriter, r *http.Request)
	AnalysisPage(w http.ResponseWriter, r *http.Request)
}

// digitalInvestigationHandler implements IDigitalInvestigationHandler.
// The prototype is client-side, so it only needs to render the page.
type digitalInvestigationHandler struct{}

// NewHandler wires the digital investigation minigame's HTTP layer.
func NewHandler() IDigitalInvestigationHandler {
	return &digitalInvestigationHandler{}
}

// Page renders the digital investigation minigame page.
func (h *digitalInvestigationHandler) Page(w http.ResponseWriter, r *http.Request) {
	pages.DigitalInvestigation().Render(r.Context(), w)
}

// AnalysisPage renders Phase 2 of the digital investigation minigame.
func (h *digitalInvestigationHandler) AnalysisPage(w http.ResponseWriter, r *http.Request) {
	pages.DigitalInvestigationAnalysis().Render(r.Context(), w)
}
