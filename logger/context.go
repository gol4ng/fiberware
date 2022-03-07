package logger

import (
	"context"
	"github.com/gol4ng/logger"
)

type ctxKey uint

const (
	loggerKey ctxKey = iota
)

func LoggerToContext(ctx context.Context, credential logger.LoggerInterface) context.Context {
	return context.WithValue(ctx, loggerKey, credential)
}

func LoggerFromContext(ctx context.Context) logger.LoggerInterface {
	if ctx == nil {
		return nil
	}
	value := ctx.Value(loggerKey)
	if value == nil {
		return nil
	}
	credential, ok := value.(logger.LoggerInterface)
	if !ok {
		return nil
	}

	return credential
}
