package controller

import (
	"github.com/gin-gonic/gin"
)

const (
	rootPageRoute = "/"
)

func Router(r *gin.Engine) {
	r.GET(rootPageRoute, AuthHandler())
}
