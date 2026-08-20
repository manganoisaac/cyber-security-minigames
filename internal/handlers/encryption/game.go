package encryption

import (
	"net/http"

	"cyber-security-minigames/internal/handlers"
)

// game adapts the encryption minigame to handlers.Game so it can be registered
// alongside the other minigames.
type game struct {
	handler IEncryptionHandler
}

// NewGame returns the encryption minigame as a handlers.Game, ready to register
// its page and htmx endpoints onto the app mux.
func NewGame(h IEncryptionHandler) handlers.Game {
	return &game{handler: h}
}

// Register mounts the encryption minigame's routes, delegating to Register.
func (g *game) Register(mux *http.ServeMux) {
	Register(mux, g.handler)
}
