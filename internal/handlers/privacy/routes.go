package privacy

import "net/http"

func Register(mux *http.ServeMux, h IPrivacyHandler) {
	mux.HandleFunc("GET /privacy", h.Page)
	mux.HandleFunc("POST /privacy/clickpost/{id}", h.ClickPost)
}
