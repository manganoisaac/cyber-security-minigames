package network_defense

import (
	"net/http"

	"cyber-security-minigames/internal/handlers"
)

// game adapts the network-defense minigame to handlers.Game so it can be registered
// alongside the other minigames.
type game struct {
	handler INetworkDefenseHandler
}

// NewGame returns the network-defense minigame as a handlers.Game, ready to register
// its page and htmx endpoints onto the app mux.
func NewGame(h INetworkDefenseHandler) handlers.Game {
	return &game{handler: h}
}

// Register mounts the network-defense minigame's routes, delegating to Register.
func (g *game) Register(mux *http.ServeMux) {
	Register(mux, g.handler)
}
