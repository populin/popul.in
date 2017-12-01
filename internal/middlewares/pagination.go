package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/populin/popul.in/internal/handlers"
	"github.com/populin/popul.in/internal/request"
)

// Pagination is a handler Gin func to extract pagination params from the request
func Pagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		p, err := request.ExtractPagination(c.Request.URL.Query())

		if decodeErrors, ok := err.(form.DecodeErrors); ok {
			b := handlers.NewErrorBuilder()
			for _, err := range decodeErrors {
				b.AddError(http.StatusBadRequest, err.Error())
			}
			handler := handlers.Error(b.Errors...)
			handler(c)
			c.Abort()
			return
		}

		c.Set("pagination", p)
	}
}
