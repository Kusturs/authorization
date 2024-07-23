package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/mtank-group/auth-go/src/config"
	"github.com/mtank-group/auth-go/src/internal/controller"
	"github.com/mtank-group/auth-go/src/pkg/logger"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	utc, err := time.LoadLocation(time.UTC.String())
	if err != nil {
		log.Fatal(fmt.Sprintf("app - Run - time.LoadLocation: %s", err.Error()))
	}

	time.Local = utc
	gin.SetMode(cfg.App.Mode)

	engine := gin.New()
	engine.Use(gin.Logger())

	//hnd := controller.New(log)

	controller.Router(engine)

	if err := engine.Run(cfg.App.Port); err != nil {
		log.Fatal(fmt.Sprintf("app - Run - engine.Run: %s", err.Error()))
	}
}
