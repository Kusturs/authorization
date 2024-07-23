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
	}

	App struct {
		Mode string `env-required:"true" env:"APP_MODE" env-upd:"true"`
		Port string `env-required:"true" env:"APP_PORT" env-upd:"true"`
	}

	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}
)

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
