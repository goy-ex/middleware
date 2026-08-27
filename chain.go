package middleware

import (
	"net/http"
	"slices"
)

func Chain(base http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range slices.Backward(middlewares) {
		base = mw(base)
	}

	return base
}
