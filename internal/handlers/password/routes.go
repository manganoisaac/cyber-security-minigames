package password

import "net/http"

// Register wires up the password page and all of its htmx endpoints,
// delegating to the given handler.
func Register(mux *http.ServeMux, h IPasswordHandler) {
	mux.HandleFunc("GET /password", h.Page)
	mux.HandleFunc("GET /password/selectitem", h.SelectItem)
	mux.HandleFunc("POST /password/guess/{guess}", h.AddGuess)
	mux.HandleFunc("POST /password/addguess/{guess}", h.AddGuess)
	mux.HandleFunc("POST /password/backspace", h.Backspace)
	mux.HandleFunc("GET /password/changeview/{view_id}", h.ChangeView)
}
