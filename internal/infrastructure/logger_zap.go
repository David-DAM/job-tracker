package infrastructure

import (
	"context"
	"job-tracker/internal/domain"

	"go.uber.org/zap"
)

type LoggerZap struct {
	log *zap.Logger
}

func NewLoggerZap(log *zap.Logger) domain.Logger {
	return &LoggerZap{log: log}
}

func (l *LoggerZap) Info(ctx context.Context, msg string, fields ...domain.Field) {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}

	l.log.With(zap.Any("context", ctx)).Info(msg, zapFields...)
}

func (l *LoggerZap) Debug(ctx context.Context, msg string, fields ...domain.Field) {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}

	l.log.With(zap.Any("context", ctx)).Debug(msg, zapFields...)
}

func (l *LoggerZap) Error(ctx context.Context, msg string, err error, fields ...domain.Field) {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}

	zapFields = append(zapFields, zap.Error(err))

	l.log.With(zap.Any("context", ctx)).Error(msg, zapFields...)
}
