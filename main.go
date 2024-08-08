package main

import (
	"app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	routes.RegisterRoutes(router)

	router.Run(":8080")

}
