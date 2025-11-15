package domain

import (
	"context"
	"timo/dto"
	"timo/models"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	LoginWithPassword(ctx context.Context, req *dto.LoginWithPasswordRequest) (*dto.LoginResponse, error)
	LoginWithGoogle(ctx context.Context, req *dto.LoginWithGoogleRequest) (*dto.LoginResponse, error)
}
