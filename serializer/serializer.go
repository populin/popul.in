package serializer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Serializer is a registry of FormatHandler structs
type Serializer struct {
	handlers []FormatHandler
}

// FormatHandler defines the behavior of a handler
type FormatHandler interface {
	Supports(format string) bool
	Handle(o interface{}) ([]byte, error)
}

// NewSerializer is the factory of Serializer
func NewSerializer() Serializer {
	s := Serializer{
		handlers: []FormatHandler{
			GeoJSONHandler{},
			JSONAPIHandler{},
		},
	}
	return s
}

// Serialize will fetch all handlers and try to find one which supports the requested format
func (s Serializer) Serialize(c *gin.Context, o interface{}) ([]byte, error) {
	for _, handler := range s.handlers {
		if handler.Supports(c.GetHeader("Accept")) {
			return handler.Handle(o)
		}
	}

	return []byte{}, errors.New("no handler found supporting the requested format")
}
