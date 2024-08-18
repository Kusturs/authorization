package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/solndev/auth-go/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT username, password FROM users WHERE username=$1`
	row := r.db.QueryRow(ctx, query, username)

	var user models.User
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, user.Username, user.Password)
	return err
}
