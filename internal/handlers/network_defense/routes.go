package network_defense

import "net/http"

// Register wires up the network-defense page, delegating to the given handler.
// The minigame runs entirely client-side, so the page is its only route.
func Register(mux *http.ServeMux, h INetworkDefenseHandler) {
	mux.HandleFunc("GET /network-defense", h.Page)
}
