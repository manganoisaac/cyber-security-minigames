package email_client

import (
	"net/http/httptest"
	"testing"

	"cyber-security-minigames/data"
)

func TestLevelToggle(t *testing.T) {
	r := httptest.NewRequest("GET", "/email-client", nil)
	w := httptest.NewRecorder()
	s := getState(w, r)

	if got := s.level(); got != 1 {
		t.Fatalf("initial level = %d, want 1", got)
	}
	s.setLevel(2)
	if got := s.level(); got != 2 {
		t.Fatalf("after setLevel, level = %d, want 2", got)
	}
}

func TestSessionRecallByCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/email-client", nil)
	w := httptest.NewRecorder()
	first := getState(w, r)
	phishing, _, _ := data.PhishingEmail()
	first.setActiveEmail(phishing.ID)

	r2 := httptest.NewRequest("POST", "/email-check/"+phishing.ID, nil)
	for _, ck := range w.Result().Cookies() {
		r2.AddCookie(ck)
	}
	w2 := httptest.NewRecorder()
	second := getState(w2, r2)
	if second != first {
		t.Fatal("cookie did not recall the same session")
	}
	if second.activeEmail() != phishing.ID {
		t.Fatalf("recalled session active = %q, want %q", second.activeEmail(), phishing.ID)
	}
}
