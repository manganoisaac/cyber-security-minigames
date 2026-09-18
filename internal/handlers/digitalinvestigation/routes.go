package digitalinvestigation

import "net/http"

// Register wires up the digital investigation pages.
func Register(mux *http.ServeMux, h IDigitalInvestigationHandler) {
	mux.HandleFunc("GET /digital-investigation", h.Page)
	mux.HandleFunc("GET /digital-investigation/analysis", h.AnalysisPage)
	mux.HandleFunc("GET /whwl/digital-investigation", h.LearnedPage)
}
