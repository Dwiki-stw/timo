package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"timo/dto"
	"timo/helper"
	"timo/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMocks func(svc *mocks.AuthServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "payload validation failed",
			body:       `{}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   helper.VALIDATION_ERROR,
		},
		{
			name: "service return error",
			body: `{"name": "user test", "email": "test@example.com", "password": "test123"}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				svc.On("Register", mock.Anything, mock.AnythingOfType("*dto.RegisterRequest")).
					Return(nil, helper.NewAppError(helper.EMAIL_EXIST, "email already registered", nil))
			},
			wantCode: http.StatusConflict,
			wantBody: helper.EMAIL_EXIST,
		},
		{
			name: "success",
			body: `{"name": "user test", "email": "test@example.com", "password": "test123"}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				svc.On("Register", mock.Anything, mock.AnythingOfType("*dto.RegisterRequest")).
					Return(&dto.RegisterResponse{
						Uid:   "UIDtest123",
						Email: "test@example.com",
						Name:  "user test",
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `"uid":"UIDtest123"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			svc := new(mocks.AuthServiceMock)
			tt.setupMocks(svc)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			h := NewAuth(svc)
			h.Register(c)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_LoginWithGoogle(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMocks func(svc *mocks.AuthServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "payload validation failed",
			body:       `{}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   helper.VALIDATION_ERROR,
		},
		{
			name: "service return error",
			body: `{"id_token": "token-googleTest"}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				svc.On("LoginWithGoogle", mock.Anything, mock.AnythingOfType("*dto.LoginWithGoogleRequest")).
					Return(nil, helper.NewAppError(helper.NOT_FOUND, "user not found", assert.AnError))
			},
			wantCode: http.StatusNotFound,
			wantBody: helper.NOT_FOUND,
		},
		{
			name: "success",
			body: `{"id_token": "token-googleTest"}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				svc.On("LoginWithGoogle", mock.Anything, mock.AnythingOfType("*dto.LoginWithGoogleRequest")).
					Return(&dto.LoginResponse{
						Uid:   "UIDtest123",
						Name:  "user test",
						Token: "tokentest123",
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `"token":"tokentest123"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			svc := new(mocks.AuthServiceMock)
			tt.setupMocks(svc)

			req := httptest.NewRequest(http.MethodPost, "/login/google", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			h := NewAuth(svc)
			h.LoginWithGoogle(c)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_LoginWithPasword(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		setupMocks func(svc *mocks.AuthServiceMock)
		wantCode   int
		wantBody   string
	}{
		{
			name:       "payload validation failed",
			body:       `{}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {},
			wantCode:   http.StatusBadRequest,
			wantBody:   helper.VALIDATION_ERROR,
		},
		{
			name: "service return error",
			body: `{"email": "test@example.com", "password": "test123"}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				svc.On("LoginWithPassword", mock.Anything, mock.AnythingOfType("*dto.LoginWithPasswordRequest")).
					Return(nil, helper.NewAppError(helper.LOGIN_ERROR, "invalid email or password", assert.AnError))
			},
			wantCode: http.StatusBadRequest,
			wantBody: helper.LOGIN_ERROR,
		},
		{
			name: "success",
			body: `{"email": "test@example.com", "password": "test123"}`,
			setupMocks: func(svc *mocks.AuthServiceMock) {
				svc.On("LoginWithPassword", mock.Anything, mock.AnythingOfType("*dto.LoginWithPasswordRequest")).
					Return(&dto.LoginResponse{
						Uid:   "UIDtest123",
						Name:  "user test",
						Token: "tokentest123",
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: `"token":"tokentest123"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			svc := new(mocks.AuthServiceMock)
			tt.setupMocks(svc)

			req := httptest.NewRequest(http.MethodPost, "/login/password", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			h := NewAuth(svc)
			h.LoginWithPassword(c)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantBody)
			svc.AssertExpectations(t)
		})
	}
}
