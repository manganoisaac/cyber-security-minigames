package encryption

import "net/http"

// Register wires up the encryption page and all of its htmx endpoints,
// delegating to the given handler.
func Register(mux *http.ServeMux, h IEncryptionHandler) {
	mux.HandleFunc("GET /encryption", h.Page)
}
