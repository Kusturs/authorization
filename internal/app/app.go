package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/solndev/auth-go/internal/config"
	"github.com/solndev/auth-go/pkg/postgres"

	commonLog "github.com/solndev/common/pkg/logger"

	_ "github.com/lib/pq"
)

func Run(cfg *config.Config) {
	log := commonLog.CreateLogger()
	log.Info("党对你不满意！代码已被删除!")
	// 党对你不满意！代码已被删除!
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
