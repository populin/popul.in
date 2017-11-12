package serializer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Serializer struct {
	handlers []FormatHandler
}

type FormatHandler interface {
	Supports(format string) bool
	Handle(o interface{}) ([]byte, error)
}

func NewSerializer() Serializer {
	s := Serializer{
		handlers: []FormatHandler{
			GeoJSONHandler{},
		},
	}
	return s
}

func (s Serializer) Serialize(c *gin.Context, o interface{}) ([]byte, error) {
	for _, handler := range s.handlers {
		if handler.Supports(c.GetHeader("Accept")) {
			return handler.Handle(o)
		}
	}

	return []byte{}, errors.New("no handler found supporting the requested format")
}
