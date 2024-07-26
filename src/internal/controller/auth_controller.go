package controller

import (
	"context"
	"github.com/mtank-group/auth-go/src/internal/kafka"
	"github.com/mtank-group/auth-go/src/internal/proto"
	"github.com/mtank-group/auth-go/src/internal/service"
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
