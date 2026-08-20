package handlers

import (
	"net/http"

	"cyber-security-minigames/views/pages"

	"github.com/a-h/templ"
)

// Game is a self-contained minigame that mounts its page and htmx routes
// onto the app mux.
type Game interface {
	Register(mux *http.ServeMux)
}

// Routes wires up the page routes plus static asset serving and the given
// games, returning the handler for the server to run. New games are added by
// passing them in from main.go; this function never changes when a game is
// added or removed.
func Routes(games ...Game) http.Handler {
	mux := http.NewServeMux()

	// Static assets (css / js) served from ./public at /public/…
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("GET /public/", http.StripPrefix("/public/", fs))

	// Core pages. "GET /{$}" matches only the exact root path, so it
	// doesn't swallow every other route.
	mux.Handle("GET /{$}", templ.Handler(pages.Landing()))
	mux.Handle("GET /chapter-select", templ.Handler(pages.ChapterSelect()))

	// Each game registers its page and htmx endpoints here.
	for _, g := range games {
		g.Register(mux)
	}

	return mux
}
