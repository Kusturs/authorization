package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/solndev/auth-go/internal/kafka"
	"github.com/solndev/auth-go/internal/proto"
	"github.com/solndev/auth-go/internal/service"
)

type AuthController struct {
	proto.UnimplementedAuthServiceServer
	userService   *service.UserService
	kafkaProducer *kafka.Producer
}

func NewAuthController(userService *service.UserService, kafkaProducer *kafka.Producer) *AuthController {
	return &AuthController{
		userService:   userService,
		kafkaProducer: kafkaProducer,
	}
}

func (c *AuthController) Authenticate(ctx context.Context, req *proto.AuthRequest) (*proto.AuthResponse, error) {
	user, err := c.userService.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		return &proto.AuthResponse{
			Success: false,
			Message: "Authentication failed",
		}, err
	}

	token, err := c.userService.GenerateJWT(user)
	if err != nil {
		return &proto.AuthResponse{
			Success: false,
			Message: "Failed to generate token",
		}, err
	}

	err = c.kafkaProducer.SendMessage("user-authenticated", user.Username)
	if err != nil {
		return &proto.AuthResponse{
			Success: false,
			Message: "Failed to send Kafka message",
		}, err
	}

	return &proto.AuthResponse{
		Token:   token,
		Success: true,
		Message: "Authenticated",
	}, nil
}

func (c *AuthController) AuthenticateHandler(ctx *gin.Context) {
	var req proto.AuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := c.Authenticate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": resp.Message})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req proto.AuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := c.userService.RegisterUser(ctx, req.Username, req.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
