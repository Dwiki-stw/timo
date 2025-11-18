package mocks

import (
	"context"
	"timo/dto"
	"timo/models"

	"github.com/stretchr/testify/mock"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func (a *AuthRepositoryMock) CreateUser(ctx context.Context, user *models.User) error {
	args := a.Called(ctx, user)
	return args.Error(0)
}

func (a *AuthRepositoryMock) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := a.Called(ctx, email)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

type AuthServiceMock struct {
	mock.Mock
}

func (a *AuthServiceMock) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := a.Called(ctx, req)
	if resp, ok := args.Get(0).(*dto.RegisterResponse); ok {
		return resp, args.Error(1)
	}

	return nil, args.Error(1)
}

func (a *AuthServiceMock) LoginWithPassword(ctx context.Context, req *dto.LoginWithPasswordRequest) (*dto.LoginResponse, error) {
	args := a.Called(ctx, req)
	if resp, ok := args.Get(0).(*dto.LoginResponse); ok {
		return resp, args.Error(1)
	}

	return nil, args.Error(1)
}

func (a *AuthServiceMock) LoginWithGoogle(ctx context.Context, req *dto.LoginWithGoogleRequest) (*dto.LoginResponse, error) {
	args := a.Called(ctx, req)
	if resp, ok := args.Get(0).(*dto.LoginResponse); ok {
		return resp, args.Error(1)
	}

	return nil, args.Error(1)
}
