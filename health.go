package middleware

import "net/http"

func Health(pattern string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Pattern == pattern && r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
			}

			next.ServeHTTP(w, r)
		})
	}
}
