package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mtank-group/auth-go/src/internal/service"
	"github.com/mtank-group/auth-go/src/pkg/logger"
	"net/http"
)

type Handler struct {
	service *service.TestService
	log     *logger.Logger
}

func New(service *service.TestService, log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

func (h *Handler) TestHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello World")
	}
}
