package digitalinvestigation

import (
	"net/http"

	"cyber-security-minigames/internal/handlers"
)

// game adapts the digital investigation minigame to handlers.Game so it can be
// registered alongside the other minigames.
type game struct {
	handler IDigitalInvestigationHandler
}

// NewGame returns the digital investigation minigame as a handlers.Game.
func NewGame(h IDigitalInvestigationHandler) handlers.Game {
	return &game{handler: h}
}

// Register mounts the digital investigation minigame's routes.
func (g *game) Register(mux *http.ServeMux) {
	Register(mux, g.handler)
}
