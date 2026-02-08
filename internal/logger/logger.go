package logger

import (
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	production, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return production
}
