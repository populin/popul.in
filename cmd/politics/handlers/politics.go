package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ByID returns a Politic by its ID
func ByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
