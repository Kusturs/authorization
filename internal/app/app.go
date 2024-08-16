package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/solndev/auth-go/config"
	"github.com/solndev/auth-go/internal/kafka"
	"github.com/solndev/auth-go/pkg/logger"
	"github.com/solndev/auth-go/pkg/postgres"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
	"github.com/solndev/auth-go/internal/controller"
	pb "github.com/solndev/auth-go/internal/proto"
	"github.com/solndev/auth-go/internal/repository"
	"github.com/solndev/auth-go/internal/service"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	utc, err := time.LoadLocation(time.UTC.String())
	if err != nil {
		log.Fatal(fmt.Sprintf("app - Run - time.LoadLocation: %s", err.Error()))
	}

	time.Local = utc
	gin.SetMode(cfg.App.Mode)

	pg, err := postgres.New(GetDbConnectionUrl(cfg))
	if err != nil {
		log.Fatal(fmt.Sprintf("app - Run - postgres.New: %s", err.Error()))
	}
	defer pg.Close()

	err = pg.Pool.Ping(context.Background())
	if err != nil {
		log.Fatal(fmt.Sprintf("app - Run - postgres.Ping: %s", err.Error()))
	}

	if err := runMigrations(pg); err != nil {
		log.Fatal(fmt.Sprintf("app - Run - runMigrations: %s", err.Error()))
	}

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

	// Initialize Gin HTTP server
	r := gin.Default()
	controller.Router(r, authController)

	go func() {
		httpPort := fmt.Sprintf("%s", cfg.App.Port) // fatal error :: when ":%s"
		if err := r.Run(httpPort); err != nil {
			log.Fatal("failed to run HTTP server: %v", zap.Error(err))
		}
	}()

	// Set up gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authController)

	grpcAddress := fmt.Sprintf("0.0.0.0:%s", cfg.GRPC.Port)
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatal("failed to listen: %v", zap.Error(err))
	}

	ctx := context.Background()
	log.Info(ctx, "gRPC server listening at %v", zap.String("address", lis.Addr().String()))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", zap.Error(err))
	}
}

func GetDbConnectionUrl(cfg *config.Config) string {
	return cfg.DB.ConnectionURL()
}

func runMigrations(pg *postgres.Postgres) error {
	file, err := os.Open("src/internal/migrations/migrations.sql")
	if err != nil {
		return fmt.Errorf("failed to open migration file: %w", err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if _, err := pg.Pool.Exec(context.Background(), query); err != nil {
			return fmt.Errorf("failed to execute migration query: %w", err)
		}
	}

	return nil
}
