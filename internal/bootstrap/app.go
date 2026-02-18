package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"job-tracker/internal/domain"
	"job-tracker/internal/infrastructure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type App struct {
	Logger      domain.Logger
	JobHandler  *infrastructure.JobHandler
	JobScrapper *infrastructure.JobScrapper
}

func NewApp(logger domain.Logger, jobHandler *infrastructure.JobHandler, jobScrapper *infrastructure.JobScrapper) *App {
	return &App{
		Logger:      logger,
		JobHandler:  jobHandler,
		JobScrapper: jobScrapper,
	}
}

func Start() error {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, lp, err := InitLogs()
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if err := lp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	tracer, err := InitTracer()
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if err := tracer.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	metrics, err := InitMetrics()
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if err := metrics.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}()

	config, err := LoadConfig()
	if err != nil {
		return err
	}
	db, err := NewDatabase(config)
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
	r.Use(gin.Logger())
	r.Use(
		otelgin.Middleware(
			config.AppName,
			otelgin.WithMeterProvider(metrics),
			otelgin.WithTracerProvider(tracer),
		),
	)
	app.JobHandler.RegisterRoutes(r)
	RegisterStatus(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: r,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := app.JobScrapper.InitScrape(ctx); err != nil {
			app.Logger.Error(ctx, "scrapper stopped", err)
		}
	}()

	go func() {
		app.Logger.Info(ctx, "server started")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error(ctx, "server error", err)
			stop()
		}
	}()

	<-ctx.Done()
	app.Logger.Info(ctx, "shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		app.Logger.Error(ctx, "server shutdown error", err)
	}

	wg.Wait()

	app.Logger.Info(ctx, "app stopped cleanly")
	return nil
}
