package handlers

import (
	"net/http"

	"github.com/a-h/templ"

	"cyber-security-minigames/views/components"
	"cyber-security-minigames/views/pages"
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
	mux.Handle("GET /email-client", templ.Handler(pages.EmailClient()))
	mux.Handle("POST /email-word/{id}", templ.Handler(components.EmailWord("replaced")))

	return mux
}
