package app

import (
	"context"
	"fmt"
	//"github.com/gin-gonic/gin"
	"github.com/mtank-group/auth-go/src/internal/kafka"
	"net"
	//"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/mtank-group/auth-go/src/config"
	"github.com/mtank-group/auth-go/src/internal/controller"
	pb "github.com/mtank-group/auth-go/src/internal/proto"
	"github.com/mtank-group/auth-go/src/internal/repository"
	"github.com/mtank-group/auth-go/src/internal/service"
	"github.com/mtank-group/auth-go/src/pkg/logger"
	"github.com/mtank-group/auth-go/src/pkg/postgres"

	_ "github.com/lib/pq"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	pg, err := postgres.New(
		GetDbConnectionUrl(cfg),
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("app - Run - postgres.New: %s", err.Error()))
	}
	defer pg.Close()

	err = pg.Pool.Ping(context.Background())
	if err != nil {
		log.Fatal(fmt.Sprintf("app - Run - postgres.Ping: %s", err.Error()))
	}

	// Initialize repositories and services
	userRepository := repository.NewUserRepository(pg.Pool)
	userService := service.NewUserService(userRepository, cfg.JWT.SecretKey)

	kafkaProducer, err := kafka.NewKafkaProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal("failed to create Kafka producer: %v", zap.Error(err))
	}
	defer func(kafkaProducer *kafka.Producer) {
		err := kafkaProducer.Close()
		if err != nil {
			log.Fatal("failed to close Kafka producer: %v", zap.Error(err))
		}
	}(kafkaProducer)
	authController := controller.NewAuthController(userService, kafkaProducer)

	// Set up gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authController)

	lis, err := net.Listen("tcp", cfg.App.Port)
	if err != nil {
		log.Fatal("failed to listen: %v", zap.Error(err))
	}

	ctx := context.Background()
	log.Info(ctx, "server listening at %v", zap.String("address", lis.Addr().String()))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", zap.Error(err))
	}
}

func GetDbConnectionUrl(cfg *config.Config) string {
	//if cfg.App.Mode != gin.TestMode {
	//	return cfg.PG.ConnectionURL()
	//}
	return cfg.DB.ConnectionURL()
	//return cfg.PG.ConnectionURLTest()
}
