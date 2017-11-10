package middlewares

import "github.com/gin-gonic/gin"

// GeoJSON is a Gin handler func to add response header for the content-type
func GeoJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/geo+json; charset=utf-8")
	}
}
