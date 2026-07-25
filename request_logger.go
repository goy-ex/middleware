package middleware

import (
	"context"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func RequestLogger(loggerFunc LoggerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := loggerFunc(r)
			r = r.WithContext(addLogger(r.Context(), logger))
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			logger.Info("request",
				zap.Int("status", ww.Status()),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

type RequestLoggerKey int

const requestLoggerKey RequestLoggerKey = 0

func addLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, requestLoggerKey, logger)
}

func GetLogger(ctx context.Context) *zap.Logger {
	return ctx.Value(requestLoggerKey).(*zap.Logger)
}
