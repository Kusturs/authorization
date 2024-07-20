package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mtank-group/auth-go/src/config"
	"github.com/mtank-group/auth-go/src/pkg/logger"
)

const (
	testPageRoute = "/test"
)

func Router(
	r *gin.Engine,
	cfg *config.Config,
	l *logger.Logger,
	hnd *Handler,
) {
	// r.Use(middlewareController.Metrics())
	r.GET(testPageRoute, hnd.TestHandler())
}
