package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/mtank-group/auth-go/src/config"
	"github.com/mtank-group/auth-go/src/internal/controller"
	"github.com/mtank-group/auth-go/src/pkg/logger"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)
	ctx := context.Background()

	utc, err := time.LoadLocation(time.UTC.String())
	if err != nil {
		log.Fatal(fmt.Sprintf("app - Run - time.LoadLocation: %s", err.Error()))
	}

	time.Local = utc
	gin.SetMode(cfg.App.Mode)

	engine := gin.New(log, cfg)
	engine.Use(gin.Logger())
	controller.Router(log)
}
