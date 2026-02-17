package domain

import (
	"context"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, err error, fields ...Field)
	Debug(ctx context.Context, msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value any
}
