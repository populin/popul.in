package handlers

import (
	"fmt"

	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/populin/popul.in/internal/constants"
)

// Error is a Gin handler func to create an error response
func Error(errs ...*jsonapi.ErrorObject) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", constants.JSONAPI))

		if len(errs) == 1 {
			firstStatus, _ := strconv.Atoi(errs[0].Status)
			c.Status(firstStatus)
		} else {
			c.Status(http.StatusBadRequest)
		}

		jsonapi.MarshalErrors(c.Writer, errs)
		c.Abort()
	}
}

// NotFound is the 404 handler
func NotFound() func(c *gin.Context) {
	b := NewErrorBuilder()
	b.AddError(http.StatusNotFound, "")
	return Error(b.Errors...)
}

// MethodNotAllowed is the 405 handler
func MethodNotAllowed() func(c *gin.Context) {
	b := NewErrorBuilder()
	b.AddError(http.StatusMethodNotAllowed, "")
	return Error(b.Errors...)
}

// ErrorBuilder is a struct carrying jsonapi errors
type ErrorBuilder struct {
	Errors []*jsonapi.ErrorObject
}

// NewErrorBuilder returns a pointer to an ErrorBuilder
func NewErrorBuilder() *ErrorBuilder {
	return &ErrorBuilder{}
}

// AddError makes it easier to add a simple jsonapi error to the builder
func (b *ErrorBuilder) AddError(status int, detail string) {
	b.Errors = append(
		b.Errors,
		&jsonapi.ErrorObject{
			Status: strconv.Itoa(status),
			Title:  http.StatusText(status),
			Detail: detail,
		},
	)
}
