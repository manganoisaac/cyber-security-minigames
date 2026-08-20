package email_client

import "net/http"

// Register wires up the email client page and all of its htmx endpoints,
// delegating to the given handler.
func Register(mux *http.ServeMux, h IEmailHandler) {
	mux.HandleFunc("GET /email-client", h.Page)
	mux.HandleFunc("GET /email/{id}", h.Select)
	mux.HandleFunc("POST /email-phrase/{emailID}/{phraseIndex}/{select}", h.Phrase)
	mux.HandleFunc("POST /email-submit/{emailID}", h.Submit)
	mux.HandleFunc("POST /email-check/{emailID}", h.Check)
	mux.HandleFunc("DELETE /email-client/modal", h.ModalClose)
}
