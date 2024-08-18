package controller

import (
	"github.com/gin-gonic/gin"
)

const (
	authRoute     = "/authenticate"
	registerRoute = "/register"
	//confirmRoute = "/confirm"
)

func Router(r *gin.Engine, authController *AuthController) {
	r.POST(authRoute, authController.AuthenticateHandler)
	r.POST(registerRoute, authController.Register)
	//r.GET(confirmRoute, authController.ConfirmHandler)
}
