package email_client

import "net/http"

// Register wires up the email client page and all of its htmx endpoints,
// delegating to the given handler.
func Register(mux *http.ServeMux, h IEmailHandler) {
	mux.HandleFunc("GET /phishing", h.Page)
	mux.HandleFunc("GET /phishing/instructions", h.Instructions)
	mux.HandleFunc("GET /email/{id}", h.Select)
	mux.HandleFunc("POST /email-phrase/{emailID}/{phraseIndex}/{select}", h.Phrase)
	mux.HandleFunc("POST /email-submit/{emailID}", h.Submit)
	mux.HandleFunc("POST /email-check/{emailID}", h.Check)
	mux.HandleFunc("DELETE /phishing/modal", h.ModalClose)

	// Keep the old address as a compatibility redirect so any stale links or
	// bookmarks still land on the phishing activity's canonical route.
	mux.HandleFunc("GET /email-client", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/phishing", http.StatusTemporaryRedirect)
	})

	// The shared legacy modal component still points here, so keep this close
	// endpoint available even though the phishing-specific modal now uses
	// /phishing/modal.
	mux.HandleFunc("DELETE /email-client/modal", h.ModalClose)
}
