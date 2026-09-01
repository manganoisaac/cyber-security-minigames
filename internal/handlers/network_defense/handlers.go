package network_defense

import (
	"net/http"

	"cyber-security-minigames/views/pages"
)

// INetworkDefenseHandler is the HTTP handler contract for the
// network-defense minigame.
type INetworkDefenseHandler interface {
	Page(w http.ResponseWriter, r *http.Request)
}

// networkDefenseHandler implements INetworkDefenseHandler. The game is
// fully client-side, so it only needs to render the page.
type networkDefenseHandler struct{}

// NewHandler wires the network-defense minigame's HTTP layer.
func NewHandler() INetworkDefenseHandler {
	return &networkDefenseHandler{}
}

// Page renders the network-defense minigame page.
func (h *networkDefenseHandler) Page(w http.ResponseWriter, r *http.Request) {
	pages.NetworkDefensePage().Render(r.Context(), w)
}
