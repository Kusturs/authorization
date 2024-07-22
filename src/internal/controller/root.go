package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RootHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Auth-go")
	}
}
