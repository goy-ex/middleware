package zap

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func ReqLogger(buildRequestLogger func(r *http.Request) *zap.Logger, skipPatterns ...string) func(http.Handler) http.Handler {
	skip := make(map[string]struct{}, len(skipPatterns))
	for _, s := range skipPatterns {
		skip[s] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := skip[r.URL.Path]
			if ok {
				next.ServeHTTP(w, r)

				return
			}

			logger := buildRequestLogger(r)

			defer func() {
				rec := recover()
				if rec == nil {
					return
				}
				logger.Error("panic", zap.Any("value", rec))
				panic(rec)
			}()

			r = r.WithContext(addReqLogger(r.Context(), logger))
			ww := wrappedResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
			start := time.Now()

			next.ServeHTTP(ww, r)
			logger.Info("request",
				zap.Int("status", ww.StatusCode),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
