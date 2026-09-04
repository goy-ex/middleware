package middleware

import (
	"context"
	"net/http"
)

type requestIDKeyType int

var requestIDKey requestIDKeyType = 0

func RequestID(genFn func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = genFn()
			}

			w.Header().Set("X-Request-ID", id)
			r = r.WithContext(context.WithValue(r.Context(), requestIDKey, id))

			next.ServeHTTP(w, r)
		})
	}
}
