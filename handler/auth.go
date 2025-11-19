package handler

import (
	"net/http"
	"timo/domain"
	"timo/dto"
	"timo/helper"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	svc domain.AuthService
}

func NewAuth(svc domain.AuthService) *Auth {
	return &Auth{svc: svc}
}

func (a *Auth) LoginWithGoogle(c *gin.Context) {
	var req dto.LoginWithGoogleRequest
	if details, err := helper.BindValidate(c, &req); err != nil {
		helper.Fail(c, http.StatusBadRequest, "payload validation failed", helper.VALIDATION_ERROR, details)
		return
	}

	resp, err := a.svc.LoginWithGoogle(c.Request.Context(), &req)
	if err != nil {
		err.(*helper.AppError).WriteError(c)
		return
	}

	helper.Ok(c, resp)
}

func (a *Auth) LoginWithPassword(c *gin.Context) {
	var req dto.LoginWithPasswordRequest
	if details, err := helper.BindValidate(c, &req); err != nil {
		helper.Fail(c, http.StatusBadRequest, "payload validation failed", helper.VALIDATION_ERROR, details)
		return
	}

	resp, err := a.svc.LoginWithPassword(c.Request.Context(), &req)
	if err != nil {
		err.(*helper.AppError).WriteError(c)
		return
	}

	helper.Ok(c, resp)
}

func (a *Auth) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if details, err := helper.BindValidate(c, &req); err != nil {
		helper.Fail(c, http.StatusBadRequest, "payload validation failed", helper.VALIDATION_ERROR, details)
		return
	}

	resp, err := a.svc.Register(c.Request.Context(), &req)
	if err != nil {
		err.(*helper.AppError).WriteError(c)
		return
	}

	helper.Ok(c, resp)
}
