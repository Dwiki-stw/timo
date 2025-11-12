package repository

import (
	"context"
	"fmt"
	"timo/domain"
	"timo/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type auth struct {
	Pool *pgxpool.Pool
}

func NewAuth(pool *pgxpool.Pool) domain.AuthRepository {
	return &auth{Pool: pool}
}

func (a *auth) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (google_id, name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, uid, created_at, updated_at
	`

	return a.Pool.QueryRow(ctx, query, user.GoogleID, user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.Uid, &user.CreatedAt, &user.UpdatedAt)
}

func (a *auth) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, uid, google_id, name, email, password_hash, created_at, updated_at
		FROM users 
		Where email = ?
	`
	var user models.User
	err := a.Pool.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Uid, &user.GoogleID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}

	return &user, err
}
