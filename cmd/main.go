package main

import (
	"log"

	"github.com/solndev/auth-go/internal/app"
	"github.com/solndev/auth-go/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app.Run(cfg)
}
