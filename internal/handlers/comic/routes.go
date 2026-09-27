package comic

import (
	"cyber-security-minigames/views/pages/comic"
	"net/http"

	"github.com/a-h/templ"
)

func Register(mux *http.ServeMux) {
	mux.Handle("GET /comic/password", templ.Handler(comic.PasswordComicPage()))
	mux.Handle("GET /comic/privacy", templ.Handler(comic.PrivacyComicPage()))
	mux.Handle("GET /comic/phishing", templ.Handler(comic.PhishingComicPage()))
	mux.Handle("GET /comic/encryption", templ.Handler(comic.EncryptionComicPage()))
	mux.Handle("GET /comic/network-defense", templ.Handler(comic.NetworkDefenceComicPage()))
	mux.Handle("GET /comic/digital-investigation", templ.Handler(comic.DigitalInvestigationComicPage()))

	mux.Handle("GET /case-closed", templ.Handler(comic.CaseClosedComicPage()))
}
