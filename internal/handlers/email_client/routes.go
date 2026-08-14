package email_client

import "net/http"

// Register wires up the email client page and all of its htmx endpoints.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /email-client", HandlePage)
	mux.HandleFunc("GET /email/{id}", HandleSelect)
	mux.HandleFunc("POST /email-phrase/{emailID}/{phraseIndex}/{select}", HandlePhrase)
	mux.HandleFunc("POST /email-submit/{emailID}", HandleSubmit)
	mux.HandleFunc("POST /email-check/{emailID}", HandleCheck)
	mux.HandleFunc("DELETE /email-client/modal", HandleModalClose)
}
