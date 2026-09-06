package middleware

import "net/http"

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			}
		}()

		next.ServeHTTP(w, r)
	})

}
