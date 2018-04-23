package middlewares

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gin-contrib/cors.v1"
)

// CORS is a handler Gin func to setup Cross-Origin Resource Sharing by allowing all origins
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	return cors.New(config)
}
