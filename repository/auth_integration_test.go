package repository

import (
	"context"
	"testing"
	"timo/helper"
	"timo/models"

	"github.com/stretchr/testify/assert"
)

func TestAuthRepo_CreateUser(t *testing.T) {
	tests := []struct {
		name string
		user *models.User
	}{
		{
			name: "create user with password",
			user: &models.User{Name: "Test User", Password: helper.Ptr("test123"), Email: "test@example.com"},
		},
		{
			name: "create user with google",
			user: &models.User{Name: "Test User", GoogleID: helper.Ptr("test123"), Email: "test@example.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := NewAuth(testDB)

			err := repo.CreateUser(ctx, tt.user)
			assert.NoError(t, err)
			assert.NotZero(t, tt.user.ID)
			assert.NotEmpty(t, tt.user.Uid)

			_, _ = testDB.Exec(ctx, "DELETE FROM users WHERE email=$1", tt.user.Email)
		})
	}
}

func TestAuthRepo_GetUserByEmail(t *testing.T) {
	ctx := context.Background()
	repo := NewAuth(testDB)

	expectedUser := &models.User{Name: "test user", Email: "test@example.com", Password: helper.Ptr("test123")}

	err := repo.CreateUser(ctx, expectedUser)
	assert.NoError(t, err)

	user, err := repo.GetUserByEmail(ctx, "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.NotZero(t, user.ID)
	assert.NotEmpty(t, user.Uid)

	_, _ = testDB.Exec(ctx, "DELETE FROM users WHERE email=$1", expectedUser.Email)
}
