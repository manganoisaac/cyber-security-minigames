package handlers

import (
	"net/http"

	"cyber-security-minigames/internal/session"
	"cyber-security-minigames/views/pages"
	"cyber-security-minigames/views/pages/comic"
	"cyber-security-minigames/views/pages/whwl"

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
	//
	// Both hub screens reset every minigame's session state before
	// rendering: on a walk-up exhibit, a new player picking a game from here
	// must never inherit progress (a level already advanced, hints already
	// collected) left behind by whoever used this browser last.
	mux.Handle("GET /{$}", resetSessionThenServe(templ.Handler(pages.CyberLanding())))
	mux.Handle("GET /chapter-select", resetSessionThenServe(templ.Handler(pages.CyberChapterSelect())))

	// Story comic pages.
	mux.Handle("GET /comic/password", templ.Handler(comic.PasswordComicPage()))
	mux.Handle("GET /comic/privacy", templ.Handler(comic.PrivacyComicPage()))
	mux.Handle("GET /comic/phishing", templ.Handler(comic.PhishingComicPage()))
	mux.Handle("GET /comic/encryption", templ.Handler(comic.EncryptionComicPage()))
	mux.Handle("GET /comic/network-defense", templ.Handler(comic.NetworkDefenceComicPage()))
	mux.Handle("GET /comic/digital-investigation", templ.Handler(comic.DigitalInvestigationComicPage()))

	// What Have We Learned pages.
	mux.Handle("GET /whwl/password", templ.Handler(whwl.PasswordLearnedPage()))
	mux.Handle("GET /whwl/privacy", templ.Handler(whwl.PrivacyLearnedPage()))
	mux.Handle("GET /whwl/phishing", templ.Handler(whwl.PhishingLearnedPage()))
	mux.Handle("GET /whwl/encryption", templ.Handler(whwl.EncryptionLearnedPage()))
	mux.Handle("GET /whwl/network-defense", templ.Handler(whwl.NetworkDefenceLearnedPage()))

	// Each game registers its page and htmx endpoints here.
	for _, g := range games {
		g.Register(mux)
	}

	return mux
}

// resetSessionThenServe wraps h so the current session's minigame state is
// cleared before h runs. Used on the two hub screens a player lands back on
// between games.
func resetSessionThenServe(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sess := session.SessionFrom(r); sess != nil {
			sess.ResetAll()
		}
		h.ServeHTTP(w, r)
	})
}
