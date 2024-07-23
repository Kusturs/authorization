package controller

import (
	"github.com/mtank-group/auth-go/src/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	log *logger.Logger
}

func New(log *logger.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func RootHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Auth-go")
	}
}
