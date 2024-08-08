package middleware

import (
	"app/utils"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AuthRequired(c *gin.Context) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.Abort()
		return
	}

	tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
	claims, err := utils.ValidateJWT(tokenString, jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("userID", claims["user_id"])
	c.Next()
}
