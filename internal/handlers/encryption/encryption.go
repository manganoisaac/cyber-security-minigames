package encryption

import (
	"net/http"

	"cyber-security-minigames/views/pages"

	"github.com/a-h/templ"
)

// Register wires the encryption minigame into the application's main router.
func Register(mux *http.ServeMux) {
	mux.Handle("GET /encryption", templ.Handler(pages.Encryption()))
}
