package session

import (
	"context"
	"net/http"
)

type ctxKey struct{}

// Middleware resolves the session for each request from the given store and
// stores it in the request context before delegating to next.
func Middleware(store IStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ctxKey{}, store.Get(w, r))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SessionFrom returns the session stored in the request context by
// Middleware, or nil if the request did not pass through it.
func SessionFrom(r *http.Request) ISession {
	sess, _ := r.Context().Value(ctxKey{}).(ISession)
	return sess
}
