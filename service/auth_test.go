package service

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"timo/dto"
	"timo/helper"
	"timo/mocks"
	"timo/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_LoginWithGoogle(t *testing.T) {
	tests := []struct {
		name       string
		validator  *mocks.MockTokenValidator
		token      *mocks.MockJwtToken
		setupMocks func(repo *mocks.AuthRepositoryMock)
		req        *dto.LoginWithGoogleRequest
		wantErr    string
	}{
		{
			name:       "user not found",
			validator:  &mocks.MockTokenValidator{Payload: nil, Err: assert.AnError},
			token:      &mocks.MockJwtToken{},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {},
			req:        &dto.LoginWithGoogleRequest{IdToken: "google_token_test"},
			wantErr:    helper.NOT_FOUND,
		},
		{
			name:      "failed to get user by email",
			validator: &mocks.MockTokenValidator{Payload: &helper.Payload{GoogleID: "123", Email: "test@example.com", Name: "user test"}, Err: nil},
			token:     &mocks.MockJwtToken{},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(nil, assert.AnError)
			},
			req:     &dto.LoginWithGoogleRequest{IdToken: "google_token_test"},
			wantErr: helper.INTERNAL_ERROR,
		},
		{
			name:      "failed to create user",
			validator: &mocks.MockTokenValidator{Payload: &helper.Payload{GoogleID: "123", Email: "test@example.com", Name: "user test"}, Err: nil},
			token:     &mocks.MockJwtToken{},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(nil, sql.ErrNoRows)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(assert.AnError)
			},
			req:     &dto.LoginWithGoogleRequest{IdToken: "google_token_test"},
			wantErr: helper.INTERNAL_ERROR,
		},
		{
			name:      "failed to create token",
			validator: &mocks.MockTokenValidator{Payload: &helper.Payload{GoogleID: "123", Email: "test@example.com", Name: "user test"}, Err: nil},
			token:     &mocks.MockJwtToken{Err: assert.AnError},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(nil, sql.ErrNoRows)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Run(func(args mock.Arguments) {
						user := args.Get(1).(*models.User)
						user.ID = 1
						user.Uid = "testUID"
						user.CreatedAt = time.Now()
						user.UpdatedAt = time.Now()
					}).
					Return(nil)
			},
			req:     &dto.LoginWithGoogleRequest{IdToken: "google_token_test"},
			wantErr: helper.INTERNAL_ERROR,
		},
		{
			name:      "success",
			validator: &mocks.MockTokenValidator{Payload: &helper.Payload{GoogleID: "123", Email: "test@example.com", Name: "user test"}, Err: nil},
			token:     &mocks.MockJwtToken{Err: nil, Token: helper.Ptr("Token test 123")},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(nil, sql.ErrNoRows)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Run(func(args mock.Arguments) {
						user := args.Get(1).(*models.User)
						user.ID = 1
						user.Uid = "testUID"
						user.CreatedAt = time.Now()
						user.UpdatedAt = time.Now()
					}).
					Return(nil)
			},
			req:     &dto.LoginWithGoogleRequest{IdToken: "google_token_test"},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.AuthRepositoryMock)
			tt.setupMocks(repo)

			svc := NewAuth(repo, helper.BcryptHasher{}, tt.validator, tt.token)
			resp, err := svc.LoginWithGoogle(context.Background(), tt.req)

			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, *tt.token.Token, resp.Token)
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
			}
		})
	}
}

func TestAuthService_LoginWithPassword(t *testing.T) {
	tests := []struct {
		name       string
		hasher     *mocks.MockHasher
		token      *mocks.MockJwtToken
		setupMocks func(repo *mocks.AuthRepositoryMock)
		req        *dto.LoginWithPasswordRequest
		wantErr    string
	}{
		{
			name:   "invalid email",
			hasher: &mocks.MockHasher{},
			token:  &mocks.MockJwtToken{},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(nil, sql.ErrNoRows)
			},
			req:     &dto.LoginWithPasswordRequest{Email: "test@example.com", Password: "passwordtest123"},
			wantErr: helper.LOGIN_ERROR,
		},
		{
			name:   "internal server error",
			hasher: &mocks.MockHasher{},
			token:  &mocks.MockJwtToken{},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(nil, assert.AnError)
			},
			req:     &dto.LoginWithPasswordRequest{Email: "test@example.com", Password: "passwordtest123"},
			wantErr: helper.INTERNAL_ERROR,
		},
		{
			name:   "invalid password",
			hasher: &mocks.MockHasher{ShouldFail: true},
			token:  &mocks.MockJwtToken{},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(&models.User{
						ID:       1,
						Uid:      "UIDtest",
						Name:     "test user",
						Email:    "test@example.com",
						Password: helper.Ptr("passwordtest123"),
					}, nil)
			},
			req:     &dto.LoginWithPasswordRequest{Email: "test@example.com", Password: "passwordtest123"},
			wantErr: helper.LOGIN_ERROR,
		},
		{
			name:   "failed to create token",
			hasher: &mocks.MockHasher{ShouldFail: false},
			token:  &mocks.MockJwtToken{Token: helper.Ptr("tokentest123"), Err: assert.AnError},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(&models.User{
						ID:       1,
						Uid:      "UIDtest",
						Name:     "test user",
						Email:    "test@example.com",
						Password: helper.Ptr("passwordtest123"),
					}, nil)
			},
			req:     &dto.LoginWithPasswordRequest{Email: "test@example.com", Password: "passwordtest123"},
			wantErr: helper.INTERNAL_ERROR,
		},
		{
			name:   "success",
			hasher: &mocks.MockHasher{ShouldFail: false},
			token:  &mocks.MockJwtToken{Token: helper.Ptr("tokentest123"), Err: nil},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(&models.User{
						ID:       1,
						Uid:      "UIDtest",
						Name:     "test user",
						Email:    "test@example.com",
						Password: helper.Ptr("passwordtest123"),
					}, nil)
			},
			req:     &dto.LoginWithPasswordRequest{Email: "test@example.com", Password: "passwordtest123"},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.AuthRepositoryMock)
			tt.setupMocks(repo)

			svc := NewAuth(repo, tt.hasher, &helper.GoogleValidator{}, tt.token)
			resp, err := svc.LoginWithPassword(context.Background(), tt.req)

			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, *tt.token.Token, resp.Token)
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name       string
		hasher     *mocks.MockHasher
		setupMocks func(repo *mocks.AuthRepositoryMock)
		req        *dto.RegisterRequest
		wantErr    string
	}{
		{
			name:   "failed to get user",
			hasher: &mocks.MockHasher{},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(nil, assert.AnError)
			},
			req:     &dto.RegisterRequest{Name: "test user", Email: "test@example.com", Password: "passwordtest123"},
			wantErr: helper.INTERNAL_ERROR,
		},
		{
			name:   "email already registered",
			hasher: &mocks.MockHasher{},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(&models.User{}, nil)
			},
			req:     &dto.RegisterRequest{Name: "test user", Email: "test@example.com", Password: "passwordtest123"},
			wantErr: helper.EMAIL_EXIST,
		},
		{
			name:   "failed to hash password",
			hasher: &mocks.MockHasher{ShouldFail: true},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(nil, sql.ErrNoRows)
			},
			req:     &dto.RegisterRequest{Name: "test user", Email: "test@example.com", Password: "passwordtest123"},
			wantErr: helper.INTERNAL_ERROR,
		},
		{
			name:   "failed to create user",
			hasher: &mocks.MockHasher{ShouldFail: false},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(nil, sql.ErrNoRows)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(assert.AnError)
			},
			req:     &dto.RegisterRequest{Name: "test user", Email: "test@example.com", Password: "passwordtest123"},
			wantErr: helper.INTERNAL_ERROR,
		},
		{
			name:   "success",
			hasher: &mocks.MockHasher{ShouldFail: false},
			setupMocks: func(repo *mocks.AuthRepositoryMock) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").
					Return(nil, sql.ErrNoRows)
				repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).
					Run(func(args mock.Arguments) {
						user := args.Get(1).(*models.User)
						user.Uid = "UIDtest"
						user.Email = "test@example.com"
						user.Name = "test user"
					}).
					Return(nil)
			},
			req:     &dto.RegisterRequest{Name: "test user", Email: "test@example.com", Password: "passwordtest123"},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.AuthRepositoryMock)
			tt.setupMocks(repo)

			svc := NewAuth(repo, tt.hasher, &helper.GoogleValidator{}, &helper.JwtToken{})
			resp, err := svc.Register(context.Background(), tt.req)

			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Email, resp.Email)
				assert.Equal(t, tt.req.Name, resp.Name)
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.wantErr, err.(*helper.AppError).Code)
			}
		})
	}
}
