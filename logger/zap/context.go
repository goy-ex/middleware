package zap

import (
	"context"

	"go.uber.org/zap"
)

type ReqLoggerKey int

const reqLoggerKey ReqLoggerKey = 0

func addReqLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, reqLoggerKey, logger)
}

func GetReqLogger(ctx context.Context) *zap.Logger {
	return ctx.Value(reqLoggerKey).(*zap.Logger)
}
