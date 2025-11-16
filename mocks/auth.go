package mocks

import (
	"context"
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
