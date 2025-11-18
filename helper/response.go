package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

func Ok[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Response[T]{Status: "success", Data: data})
}

func Fail(c *gin.Context, httpCode int, msg, code string, details any) {
	c.JSON(httpCode, Response[struct{}]{
		Status:  "error",
		Message: msg,
		Error:   &Error{Code: code, Details: details},
	})
}
