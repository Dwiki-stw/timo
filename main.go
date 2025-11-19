package main

import (
	"fmt"
	"log"
	"timo/config"
	"timo/database"
	"timo/handler"
	"timo/helper"
	"timo/repository"
	"timo/routes"
	"timo/service"

	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.Load()
	pool := database.GetConnection(conf.DB)

	//repo
	authRepo := repository.NewAuth(pool)

	//service
	authSvc := service.NewAuth(authRepo, helper.BcryptHasher{}, helper.NewGoogleValidator(""), helper.NewJwtToken(conf.JwtKey))

	//handler
	authH := handler.NewAuth(authSvc)

	handlers := &routes.Handlers{
		AuthHandler: *authH,
	}

	r := gin.Default()
	routes.SetupRoutes(r, handlers)

	r.SetTrustedProxies(nil)

	addresss := fmt.Sprintf("%s:%s", conf.App.Host, conf.App.Port)
	if err := r.Run(addresss); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
