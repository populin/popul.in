package handlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is a JSON response for errors
type ErrorResponse struct {
	Message string `json:"message"`
}

// ValidationErrorResponse is a JSON response for validation errors
type ValidationErrorResponse struct {
	ErrorResponse
	Errors map[string]string `json:"errors"`
}

// Error is a Gin handler func to create an error response
func Error(status int, err error) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.AbortWithStatusJSON(status, ErrorResponse{err.Error()})
	}
}

// ValidationError is a Gin handler func to create a validation-related error response
func ValidationError(status int, errors map[string]error) func(c *gin.Context) {
	return func(c *gin.Context) {

		keys := make([]string, 0, len(errors))
		messages := make(map[string]string)
		for k, err := range errors {
			keys = append(keys, k)
			messages[k] = err.Error()
		}

		message := fmt.Sprintf("Invalid values provided for %s", strings.Join(keys, ", "))

		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.AbortWithStatusJSON(status, ValidationErrorResponse{
			ErrorResponse: ErrorResponse{
				Message: message,
			},
			Errors: messages,
		})
	}
}
