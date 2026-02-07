package bootstrap

import (
	"fmt"
	"job-tracker/internal/config"
	"job-tracker/internal/domain"

	"github.com/gin-gonic/gin"
)

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

	jobHandler := InitJobHandler(db)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	jobHandler.RegisterRoutes(r)

	addr := fmt.Sprintf(":%d", cfg.Port)
	return r.Run(addr)
}
