package zap

import (
	"context"
	"net/http"
	"strings"
	"time"

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

type loggerKeyType int

const loggerKey loggerKeyType = iota

func RequestLogger(buildReqLogger func(r *http.Request) *zap.Logger, skipPatterns ...string) func(http.Handler) http.Handler {
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
			logger.Info("request received")

			r = r.WithContext(context.WithValue(r.Context(), loggerKey, logger))
			ww := &wrappedResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
			start := time.Now()

			defer func() {
				if rec := recover(); rec != nil {
					logger.Error(
						"panic",
						zap.Any("value", rec),
						zap.Int("status", http.StatusInternalServerError),
						zap.Duration("duration", time.Since(start)),
					)
					panic(rec)
				}
			}()

			next.ServeHTTP(ww, r)
			logger.Info("request processed",
				zap.Int("status", ww.StatusCode),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

func LoggerFrom(ctx context.Context) (*zap.Logger, bool) {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	return logger, ok
}
