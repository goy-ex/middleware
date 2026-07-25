package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

type LoggerFunc func(r *http.Request) *zap.Logger
