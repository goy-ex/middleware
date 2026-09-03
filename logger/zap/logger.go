package zap

import (
	"net/http"
	"strings"
	"time"

	ctxzap "github.com/goy-ex/middleware/logger/zap/context"
	"go.uber.org/zap"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func ReqLogger(buildReqLogger func(r *http.Request) *zap.Logger, skipPatterns ...string) func(http.Handler) http.Handler {
	skip := make(map[string]struct{}, len(skipPatterns))
	for _, s := range skipPatterns {
		skip[strings.TrimLeft(s, "/")] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := skip[strings.TrimLeft(r.URL.Path, "/")]
			if ok {
				next.ServeHTTP(w, r)

				return
			}

			logger := buildReqLogger(r)

			r = r.WithContext(ctxzap.WithLogger(r.Context(), logger))
			ww := &wrappedResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
			start := time.Now()

			defer func() {
				rec := recover()
				if rec == nil {
					return
				}
				logger.Error("panic", zap.Any("value", rec), zap.Duration("duration", time.Since(start)))
				panic(rec)
			}()

			next.ServeHTTP(ww, r)
			logger.Info("request",
				zap.Int("status", ww.StatusCode),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
