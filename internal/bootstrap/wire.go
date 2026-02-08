//go:build wireinject
// +build wireinject

package bootstrap

import (
	"job-tracker/internal/application"
	"job-tracker/internal/infrastructure"
	"job-tracker/internal/logger"

	"github.com/google/wire"
	"gorm.io/gorm"
)

//go:generate wire
func InitJobHandler(db *gorm.DB) *infrastructure.JobHandler {
	wire.Build(
		logger.NewLogger,
		infrastructure.NewJobRepository,
		application.NewJobService,
		infrastructure.NewJobHandler,
	)
	return &infrastructure.JobHandler{}
}
