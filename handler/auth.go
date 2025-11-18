package handler

import (
	"net/http"
	"timo/domain"
	"timo/dto"
	"timo/helper"

	"github.com/gin-gonic/gin"
)

type auth struct {
	svc domain.AuthService
}

func NewAuth(svc domain.AuthService) *auth {
	return &auth{svc: svc}
}

func (a *auth) LoginWithGoogle(c *gin.Context) {
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

func (a *auth) LoginWithPassword(c *gin.Context) {
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

func (a *auth) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if details, err := helper.BindValidate(c, &req); err != nil {
		helper.Fail(c, http.StatusBadRequest, "payload validation failed", helper.VALIDATION_ERROR, details)
		return
	}

	resp, err := a.svc.Register(c.Request.Context(), &req)
	if err != nil {
		err.(*helper.AppError).WriteError(c)
	}

	helper.Ok(c, resp)
}
