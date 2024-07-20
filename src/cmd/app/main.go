package main

import (
	"log"

	"github.com/mtank-group/auth-go/src/config"
	"github.com/mtank-group/auth-go/src/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
