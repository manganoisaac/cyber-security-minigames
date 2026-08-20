package handlers

import (
	"net/http"

	"cyber-security-minigames/internal/handlers/email_client"
	"cyber-security-minigames/internal/handlers/encryption"
	"cyber-security-minigames/views/pages"

	"github.com/a-h/templ"
)

// Routes wires up all page routes plus static asset serving and returns
// the handler for the server to run. Keeping this here (rather than in
// main.go) keeps the entry point small and gives the team one place to
// register new pages.
func Routes() http.Handler {
	mux := http.NewServeMux()

	// Static assets (css / js) served from ./public at /public/…
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("GET /public/", http.StripPrefix("/public/", fs))

	// Core pages. "GET /{$}" matches only the exact root path, so it
	// doesn't swallow every other route.
	mux.Handle("GET /{$}", templ.Handler(pages.Landing()))
	mux.Handle("GET /chapter-select", templ.Handler(pages.ChapterSelect()))

	// Game routes. Each game registers its page and htmx endpoints here.
	email_client.Register(mux)
	encryption.Register(mux)

	return mux
}
