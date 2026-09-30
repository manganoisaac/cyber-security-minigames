package privacy

import "net/http"

func Register(mux *http.ServeMux, h IPrivacyHandler) {
	mux.HandleFunc("GET /privacy", h.Page)
	mux.HandleFunc("GET /privacy/instructions", h.Instructions)
	mux.HandleFunc("POST /privacy/clickpost/{id}", h.ClickPost)

	mux.HandleFunc("DELETE /privacy/modal", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
