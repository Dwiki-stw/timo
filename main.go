package main

import (
	"fmt"
	"log"
	"timo/config"
	"timo/database"
	"timo/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.Load()
	_ = database.GetConnection(conf.DB)

	r := gin.Default()
	routes.SetupRoutes(r)

	r.SetTrustedProxies(nil)

	addresss := fmt.Sprintf("%s:%s", conf.App.Host, conf.App.Port)
	if err := r.Run(addresss); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
