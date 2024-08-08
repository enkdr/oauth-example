package routes

import (
	"app/handlers"
	"app/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	auth := router.Group("/auth")
	{
		auth.GET("/github", handlers.GitHubLogin)
		auth.GET("/github/callback", handlers.GitHubCallback)
		// auth.GET("/google", handlers.GoogleLogin)
		// auth.GET("/google/callback", handlers.GoogleCallback)
	}

	user := router.Group("/user")
	user.Use(middleware.AuthRequired)
	{
		user.GET("/", handlers.UserHome)
	}

	router.GET("/", handlers.Home)
}
