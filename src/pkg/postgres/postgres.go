// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMinPoolSize  = 1
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
	_idleTimeoutMinutes  = time.Minute * 1
)

// Postgres -.
type Postgres struct {
	Builder            squirrel.StatementBuilderType
	Pool               *pgxpool.Pool
	connTimeout        time.Duration
	idleTimeoutMinutes time.Duration
	maxPoolSize        int
	minPoolSize        int
	connAttempts       int
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		minPoolSize:        _defaultMinPoolSize,
		maxPoolSize:        _defaultMaxPoolSize,
		connAttempts:       _defaultConnAttempts,
		connTimeout:        _defaultConnTimeout,
		idleTimeoutMinutes: _idleTimeoutMinutes,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MinConns = int32(pg.minPoolSize)
	poolConfig.MaxConns = int32(pg.maxPoolSize)
	poolConfig.MaxConnIdleTime = pg.idleTimeoutMinutes

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *Postgres) Ping(ctx context.Context) error {
	err := p.Pool.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}

func CreateDatabaseIfNotExists(url, dbName string) error {
	// Подключение без указания базы данных
	pg, err := New(url)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer pg.Close()

	// Проверка существования базы данных
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%s'));", dbName)
	err = pg.Pool.QueryRow(context.Background(), query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	// Создание базы данных, если она не существует
	if !exists {
		_, err = pg.Pool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}
