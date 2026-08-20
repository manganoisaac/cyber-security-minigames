package email_client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cyber-security-minigames/data"
	"cyber-security-minigames/internal/session"
)

func TestLevelToggle(t *testing.T) {
	r := httptest.NewRequest("GET", "/email-client", nil)
	w := httptest.NewRecorder()

	var st IEmailGameState
	middleware := session.Middleware(session.Store())
	middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		st = stateForRequest(r)
	})).ServeHTTP(w, r)

	if got := st.Level(); got != 1 {
		t.Fatalf("initial level = %d, want 1", got)
	}
	st.SetLevel(2)
	if got := st.Level(); got != 2 {
		t.Fatalf("after SetLevel, level = %d, want 2", got)
	}
}

func TestSessionRecallByCookie(t *testing.T) {
	middleware := session.Middleware(session.Store())

	r := httptest.NewRequest("GET", "/email-client", nil)
	w := httptest.NewRecorder()

	var first IEmailGameState
	middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		first = stateForRequest(r)
	})).ServeHTTP(w, r)
	phishing, _, _ := data.PhishingEmail()
	first.SetActiveEmail(phishing.ID)

	r2 := httptest.NewRequest("POST", "/email-check/"+phishing.ID, nil)
	for _, ck := range w.Result().Cookies() {
		r2.AddCookie(ck)
	}
	w2 := httptest.NewRecorder()

	var second IEmailGameState
	middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		second = stateForRequest(r)
	})).ServeHTTP(w2, r2)

	if second != first {
		t.Fatal("cookie did not recall the same session")
	}
	if second.ActiveEmail() != phishing.ID {
		t.Fatalf("recalled session active = %q, want %q", second.ActiveEmail(), phishing.ID)
	}
}
