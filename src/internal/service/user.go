package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mtank-group/auth-go/src/internal/models"
	"github.com/mtank-group/auth-go/src/internal/repository"
)

type UserService struct {
	userRepo  *repository.UserRepository
	secretKey string
}

func NewUserService(userRepo *repository.UserRepository, secretKey string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		secretKey: secretKey,
	}
}

func (s *UserService) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *UserService) GenerateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(s.secretKey))
}

func (s *UserService) RegisterUser(ctx context.Context, username, password string) error {
	return s.userRepo.CreateUser(ctx, &models.User{Username: username, Password: password}) // Передайте ctx здесь
}
