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
func InitApp(db *gorm.DB) *App {
	wire.Build(
		logger.NewLogger,
		infrastructure.NewJobRepository,
		application.NewJobService,
		infrastructure.NewJobHandler,
		wire.Struct(new(App), "*"),
	)
	return nil
}
