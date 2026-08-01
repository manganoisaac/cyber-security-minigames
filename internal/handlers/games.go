package handlers

import (
	"net/http"

	"github.com/a-h/templ"

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
	mux.Handle("GET /activities", templ.Handler(pages.ActivitiesIndex()))
	mux.Handle("GET /credits", templ.Handler(pages.Credits()))

	// One route per activity, generated from the shared Activities list
	// so the menu and the routes never drift apart.
	for _, a := range pages.Activities {
		meta := pages.ActivityMeta{Number: a.Number, Title: a.Title, Topic: a.Topic}
		mux.Handle("GET /activity/"+a.Slug, templ.Handler(pages.Activity(meta)))
	}

	return mux
}
