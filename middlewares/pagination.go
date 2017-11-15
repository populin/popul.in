package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/populin/popul.in/handlers"
	"github.com/populin/popul.in/request"
)

// Pagination is a handler Gin func to extract pagination params from the request
func Pagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, err := request.ExtractPagination(c.Request.URL.Query())

		if decodeErrors, ok := err.(form.DecodeErrors); ok {
			handlers.ValidationError(http.StatusBadRequest, decodeErrors)
			c.Abort()
			return
		}

		c.Set("pagination", p)
	}
}
