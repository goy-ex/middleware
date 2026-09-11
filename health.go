package middleware

import (
	"net/http"
	"strings"
)

func Health(path string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.TrimLeft(r.URL.Path, "/") == strings.TrimLeft(path, "/") {
				w.WriteHeader(http.StatusOK)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
