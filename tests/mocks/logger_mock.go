package mocks

import (
	"context"
	"job-tracker/internal/domain"
)

type LoggerMock struct{}

func (l *LoggerMock) Info(ctx context.Context, msg string, fields ...domain.Field) {

}

func (l *LoggerMock) Debug(ctx context.Context, msg string, fields ...domain.Field) {

}

func (l *LoggerMock) Error(ctx context.Context, msg string, err error, fields ...domain.Field) {

}
