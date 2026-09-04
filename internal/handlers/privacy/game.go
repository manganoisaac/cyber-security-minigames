package privacy

import (
	"cyber-security-minigames/internal/handlers"
	"net/http"
)

type game struct {
	handler IPrivacyHandler
}

func NewGame(h IPrivacyHandler) handlers.Game {
	return &game{handler: h}
}

func (g *game) Register(mux *http.ServeMux) {
	Register(mux, g.handler)
}
