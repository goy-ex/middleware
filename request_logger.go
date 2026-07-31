package middleware

import (
	"context"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func RequestLogger(loggerFunc LoggerFunc, skipPaths ...string) func(http.Handler) http.Handler {
	skipMap := make(map[string]struct{}, len(skipPaths))
	for _, s := range skipPaths {
		skipMap[s] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := skipMap[r.URL.Path]
			if ok {
				next.ServeHTTP(w, r)

				return
			}

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
