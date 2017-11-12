package middlewares

import (
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/populin/popul.in/handlers"
)

// Negotiate sets the request format into the context depending on Accept Header
func Negotiate(formats ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, format := range formats {
			if format == c.GetHeader("Accept") {
				c.Header("Content-Type", fmt.Sprintf("%s; charset=utf-8", format))
				c.Next()
				return
			}
		}

		handlers.Error(http.StatusNotAcceptable, fmt.Errorf("format %s is not supported", c.GetHeader("Accept")))(c)
		c.Abort()
		return
	}
}
