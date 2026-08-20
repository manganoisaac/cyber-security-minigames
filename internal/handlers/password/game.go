package password

import (
	"net/http"

	"cyber-security-minigames/internal/handlers"
)

// game adapts the password minigame to handlers.Game so it can be registered
// alongside the other minigames.
type game struct {
	handler IPasswordHandler
}

// NewGame returns the password minigame as a handlers.Game, ready to register
// its page and htmx endpoints onto the app mux.
func NewGame(h IPasswordHandler) handlers.Game {
	return &game{handler: h}
}

// Register mounts the password minigame's routes, delegating to Register.
func (g *game) Register(mux *http.ServeMux) {
	Register(mux, g.handler)
}
