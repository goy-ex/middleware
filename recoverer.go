package middleware

import "net/http"

func Recoverer(loggerFunc LoggerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := loggerFunc(r)

			defer func() {
				rec := recover()
				if rec == nil {
					return
				}

				logger.Error("panic")
			}()

			next.ServeHTTP(w, r)
		})
	}
}
