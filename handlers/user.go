package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func UserHome(c *gin.Context) {
    userID := c.MustGet("userID").(string)
    c.HTML(http.StatusOK, "user.html", gin.H{"userID": userID})
}
