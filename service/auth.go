package service

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"timo/domain"
	"timo/dto"
	"timo/helper"
	"timo/models"
)

type auth struct {
	repo      domain.AuthRepository
	jwtKey    []byte
	hasher    helper.PasswordHasher
	validator helper.TokenValidator
}

func NewAuth(repo domain.AuthRepository) domain.AuthService {
	return &auth{repo: repo}
}

func (a *auth) LoginWithGoogle(ctx context.Context, req *dto.LoginWithGoogleRequest) (*dto.LoginResponse, error) {
	payload, err := a.validator.Validate(ctx, req.IdToken)
	if err != nil {
		return nil, helper.NewAppError(helper.NOT_FOUND, "user not found", err)
	}

	user, err := a.repo.GetUserByEmail(ctx, payload.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, helper.NewAppError(helper.INTERNAL_ERROR, "failed to get user", err)
	}

	if user == nil {
		user = &models.User{Name: payload.Name, Email: payload.Email, GoogleID: &payload.GoogleID}
		err := a.repo.CreateUser(ctx, user)
		if err != nil {
			return nil, helper.NewAppError(helper.INTERNAL_ERROR, "failed to create user", err)
		}
	}

	tokenInfo := &helper.Claims{
		UserUID: user.Uid,
		Name:    user.Name,
		Email:   user.Email,
		Exp:     time.Now().Add(time.Hour * 72).Unix(),
	}
	token, err := helper.CreateToken(tokenInfo, a.jwtKey)
	if err != nil {
		return nil, helper.NewAppError(helper.INTERNAL_ERROR, "failed to create token", err)
	}

	return &dto.LoginResponse{
		Uid:   user.Uid,
		Name:  user.Name,
		Token: *token,
	}, nil
}

func (a *auth) LoginWithPassword(ctx context.Context, req *dto.LoginWithPasswordRequest) (*dto.LoginResponse, error) {
	user, err := a.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.NewAppError(helper.LOGIN_ERROR, "invalid email or password", err)
		}
		return nil, helper.NewAppError(helper.INTERNAL_ERROR, "internal server error", err)
	}

	err = a.hasher.Compare(*user.Password, req.Password)
	if err != nil {
		return nil, helper.NewAppError(helper.LOGIN_ERROR, "invalid email or password", err)
	}

	tokenInfo := &helper.Claims{
		UserUID: user.Uid,
		Name:    user.Name,
		Email:   user.Email,
		Exp:     time.Now().Add(time.Hour * 72).Unix(),
	}

	token, err := helper.CreateToken(tokenInfo, a.jwtKey)
	if err != nil {
		return nil, helper.NewAppError(helper.INTERNAL_ERROR, "failed to create token", err)
	}

	return &dto.LoginResponse{
		Uid:   user.Uid,
		Name:  user.Name,
		Token: *token,
	}, nil
}

func (a *auth) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	existing, err := a.repo.GetUserByEmail(ctx, req.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, helper.NewAppError(helper.INTERNAL_ERROR, "internal server error", err)
	}

	if existing != nil {
		return nil, helper.NewAppError(helper.EMAIL_EXIST, "email already registered", nil)
	}

	user := &models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: &req.Password,
	}

	err = a.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, helper.NewAppError(helper.INTERNAL_ERROR, "failed to create user", err)
	}

	return &dto.RegisterResponse{
		Uid:   user.Uid,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
