package comic

import (
	"cyber-security-minigames/views/pages/comic"
	"net/http"

	"github.com/a-h/templ"
)

func Register(mux *http.ServeMux) {
	mux.Handle("GET /comic/password", templ.Handler(comic.Comic(1, "/comic/password-2")))
	mux.Handle("GET /comic/password-2", templ.Handler(comic.Comic(2, "/comic/password-3")))
	mux.Handle("GET /comic/password-3", templ.Handler(comic.Comic(3, "/comic/password-4")))
	mux.Handle("GET /comic/password-4", templ.Handler(comic.Comic(4, "/comic/password-5")))
	mux.Handle("GET /comic/password-5", templ.Handler(comic.Comic(5, "/password")))
	mux.Handle("GET /comic/privacy", templ.Handler(comic.Comic(6, "/comic/privacy-2")))
	mux.Handle("GET /comic/privacy-2", templ.Handler(comic.Comic(7, "/comic/privacy-3")))
	mux.Handle("GET /comic/privacy-3", templ.Handler(comic.Comic(8, "/privacy")))
	mux.Handle("GET /comic/phishing", templ.Handler(comic.Comic(9, "/comic/phishing-2")))
	mux.Handle("GET /comic/phishing-2", templ.Handler(comic.Comic(10, "/comic/phishing-3")))
	mux.Handle("GET /comic/phishing-3", templ.Handler(comic.Comic(11, "/email-client")))
	mux.Handle("GET /comic/encryption", templ.Handler(comic.Comic(12, "/comic/encryption-2")))
	mux.Handle("GET /comic/encryption-2", templ.Handler(comic.Comic(13, "/comic/encryption-3")))
	mux.Handle("GET /comic/encryption-3", templ.Handler(comic.Comic(14, "/comic/encryption-4")))
	mux.Handle("GET /comic/encryption-4", templ.Handler(comic.Comic(15, "/comic/encryption-5")))
	mux.Handle("GET /comic/encryption-5", templ.Handler(comic.Comic(16, "/encryption")))
	mux.Handle("GET /comic/network-defense", templ.Handler(comic.Comic(17, "/comic/network-defense-2")))
	mux.Handle("GET /comic/network-defense-2", templ.Handler(comic.Comic(18, "/comic/network-defense-3")))
	mux.Handle("GET /comic/network-defense-3", templ.Handler(comic.Comic(19, "/comic/network-defense-4")))
	mux.Handle("GET /comic/network-defense-4", templ.Handler(comic.Comic(20, "/network-defense")))
	mux.Handle("GET /comic/digital-investigation", templ.Handler(comic.Comic(21, "/comic/digital-investigation-2")))
	mux.Handle("GET /comic/digital-investigation-2", templ.Handler(comic.Comic(22, "/comic/digital-investigation-3")))
	mux.Handle("GET /comic/digital-investigation-3", templ.Handler(comic.Comic(23, "/digital-investigation")))
	mux.Handle("GET /case-closed", templ.Handler(comic.Comic(24, "/case-closed-2")))
	mux.Handle("GET /case-closed-2", templ.Handler(comic.Comic(25, "/case-closed-3")))
	mux.Handle("GET /case-closed-3", templ.Handler(comic.Comic(26, "/case-closed-4")))
	// Last panel of the game. Continues into the credits roll rather than
	// straight back to the landing screen, so finishing the whole game is the
	// only way play arrives at the credits.
	mux.Handle("GET /case-closed-4", templ.Handler(comic.Comic(27, "/credits")))

}
