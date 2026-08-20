package email_client

import (
	"net/http"

	"cyber-security-minigames/internal/handlers"
)

// game adapts the email minigame to handlers.Game so it can be registered
// alongside the other minigames.
type game struct {
	handler IEmailHandler
}

// NewGame returns the email minigame as a handlers.Game, ready to register
// its page and htmx endpoints onto the app mux.
func NewGame(h IEmailHandler) handlers.Game {
	return &game{handler: h}
}

// Register mounts the email minigame's routes, delegating to Register.
func (g *game) Register(mux *http.ServeMux) {
	Register(mux, g.handler)
}
