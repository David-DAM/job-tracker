//go:build wireinject
// +build wireinject

package bootstrap

import (
	"job-tracker/internal/application"
	"job-tracker/internal/infrastructure"

	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate wire
func InitApp(logger *zap.Logger, db *gorm.DB) *App {
	wire.Build(
		infrastructure.NewLoggerZap,
		infrastructure.NewJobRepository,
		infrastructure.NewJobScrapper,
		application.NewJobService,
		infrastructure.NewJobHandler,
		NewApp,
	)
	return nil
}
