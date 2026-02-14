//go:build wireinject
// +build wireinject

package bootstrap

import (
	"job-tracker/internal/application"
	"job-tracker/internal/domain"
	"job-tracker/internal/infrastructure"

	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:generate wire
func InitApp(logger *zap.Logger, db *gorm.DB) *App {
	wire.Build(
		domain.NewLogger,
		infrastructure.NewJobRepository,
		application.NewJobService,
		infrastructure.NewJobHandler,
		wire.Struct(new(App), "*"),
	)
	return nil
}
