package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App
		Log
		Kafka
		JWT
		DB
	}

	App struct {
		Mode string `env-required:"true" env:"APP_MODE" env-upd:"true"`
		Port string `env-required:"true" env:"APP_PORT" env-upd:"true"`
	}

	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	Kafka struct {
		Brokers []string `env-required:"true" env:"KAFKA_BROKERS"`
		Topic   string   `env-required:"true" env:"KAFKA_TOPIC"`
	}

	JWT struct {
		SecretKey string `env-required:"true" env:"JWT_SECRET_KEY"`
	}

	DB struct {
		Host     string `env-required:"true" env:"DB_HOST"`
		Port     string `env-required:"true" env:"DB_PORT"`
		User     string `env-required:"true" env:"DB_USER"`
		Password string `env-required:"true" env:"DB_PASSWORD"`
		DBName   string `env-required:"true" env:"DB_NAME"`
		SSLMode  string `env-required:"true" env:"DB_SSLMODE"`
	}
)

// NewConfig returns the application configuration.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		fmt.Printf("Environment variable error: %s, trying to read from .env", err.Error())
		err = cleanenv.ReadConfig(".env", cfg)
		if err != nil {
			if os.IsNotExist(err) {
				// when working directory is src/internal/tests
				err = cleanenv.ReadConfig(filepath.Join("..", "..", "..", ".env"), cfg)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	flag.Parse()

	return cfg, nil
}

func (c DB) ConnectionURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}
