package bootstrap

import (
	"fmt"
	"job-tracker/internal/config"
	"job-tracker/internal/domain"
	"job-tracker/internal/infrastructure"
	"job-tracker/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Logger     *zap.Logger
	JobHandler *infrastructure.JobHandler
}

func Start() error {

	cfg := config.LoadConfig()

	db, err := NewDatabase(cfg)
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&domain.Job{})
	if err != nil {
		return err
	}

	app := InitApp(db)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.TraceIDMiddleware(app.Logger))

	app.JobHandler.RegisterRoutes(r)

	addr := fmt.Sprintf(":%d", cfg.Port)
	return r.Run(addr)
}
