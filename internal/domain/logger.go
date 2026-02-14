package domain

import (
	"context"

	"go.uber.org/zap"
)

type Logger struct {
	log *zap.Logger
}

func NewLogger(log *zap.Logger) *Logger {
	return &Logger{log: log}
}

func (l *Logger) Ctx(ctx context.Context) *zap.Logger {
	return l.log.With(zap.Any("context", ctx))
}
