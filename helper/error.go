package helper

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) WriteError(c *gin.Context) {
	status := http.StatusInternalServerError
	switch e.Code {
	case INTERNAL_ERROR:
		status = http.StatusInternalServerError
	case VALIDATION_ERROR:
		status = http.StatusBadRequest
	case NOT_FOUND:
		status = http.StatusNotFound
	case LOGIN_ERROR:
		status = http.StatusBadRequest
	case EMAIL_EXIST:
		status = http.StatusConflict
	}

	Fail(c, status, e.Message, e.Code, nil)
}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

const (
	INTERNAL_ERROR   string = "INTERNAL_SERVER_ERROR"
	VALIDATION_ERROR string = "VALIDATION_ERROR"
	NOT_FOUND        string = "NOT_FOUND"
	LOGIN_ERROR      string = "LOGIN_ERROR"
	EMAIL_EXIST      string = "EMAIL_EXIST"
)
