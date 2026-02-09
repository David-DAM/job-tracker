package bootstrap

import (
	"context"
	"fmt"
	"job-tracker/internal/domain"
	"job-tracker/internal/infrastructure"
	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

type App struct {
	Logger     *zap.Logger
	JobHandler *infrastructure.JobHandler
}

func Start() error {

	cfg := LoadConfig()

	tracer, err := InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer tracer.Shutdown(context.Background())

	metrics, err := InitMetrics()
	if err != nil {
		log.Fatal(err)
	}
	defer metrics.Shutdown(context.Background())

	logger, lp, err := InitLogs()
	if err != nil {
		log.Fatal(err)
	}
	defer lp.Shutdown(context.Background())

	db, err := NewDatabase(cfg)
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&domain.Job{})
	if err != nil {
		return err
	}

	app := InitApp(logger, db)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(
		otelgin.Middleware(
			cfg.AppName,
			otelgin.WithMeterProvider(metrics),
			otelgin.WithTracerProvider(tracer),
		),
	)

	app.JobHandler.RegisterRoutes(r)

	addr := fmt.Sprintf(":%d", cfg.Port)
	return r.Run(addr)
}
