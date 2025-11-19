package routes

import (
	"net/http"
	"timo/handler"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AuthHandler handler.Auth
}

func SetupRoutes(r *gin.Engine, handlers *Handlers) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})

	r.POST("/register", handlers.AuthHandler.Register)
	r.POST("/login/password", handlers.AuthHandler.LoginWithPassword)
	r.POST("/login/google", handlers.AuthHandler.LoginWithGoogle)
}
