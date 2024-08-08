package handlers

import (
	"app/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubOauthConfig *oauth2.Config
	jwtSecret         []byte
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Initialize OAuth config using environment variables
	githubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// Initialize JWT secret
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

func GitHubLogin(c *gin.Context) {
	url := githubOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func GitHubCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "state" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	code := c.Query("code")
	token, err := githubOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	client := githubOauthConfig.Client(c, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	var userInfo map[string]interface{}
	json.Unmarshal(data, &userInfo)

	githubID := userInfo["id"].(float64)

	jwtToken, err := utils.GenerateJWT(fmt.Sprintf("%.0f", githubID), jwtSecret, time.Hour*72)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return token to the user (could also set it in a cookie)
	c.JSON(http.StatusOK, gin.H{"token": jwtToken})

}
