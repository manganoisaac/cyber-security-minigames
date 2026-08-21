package network_defense

import "net/http"

// Register wires up the network-defense page and all of its htmx endpoints,
// delegating to the given handler.
func Register(mux *http.ServeMux, h INetworkDefenseHandler) {
	mux.HandleFunc("GET /network-defense", h.Page)
}
