package encryption

import (
	"net/http"

	"cyber-security-minigames/views/pages"
)

// IEncryptionHandler is the HTTP handler contract for the encryption minigame.
type IEncryptionHandler interface {
	Page(w http.ResponseWriter, r *http.Request)
}

// encryptionHandler implements IEncryptionHandler. The game is fully
// client-side, so it only needs to render the page.
type encryptionHandler struct{}

// NewHandler wires the encryption minigame's HTTP layer.
func NewHandler() IEncryptionHandler {
	return &encryptionHandler{}
}

// Page renders the encryption minigame page.
func (h *encryptionHandler) Page(w http.ResponseWriter, r *http.Request) {
	pages.Encryption().Render(r.Context(), w)
}
