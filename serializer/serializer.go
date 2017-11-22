package serializer

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/populin/popul.in/jsonapi"
	"github.com/populin/popul.in/request"
)

// Serializer is a registry of FormatHandler structs
type Serializer struct {
	handlers []FormatHandler
}

// FormatHandler defines the behavior of a handler
type FormatHandler interface {
	Supports(format string) bool
	Handle(c *gin.Context, o interface{}) (interface{}, error)
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
			p, e := handler.Handle(c, o)

			if e != nil {
				return nil, e
			}

			if payload, ok := p.(jsonapi.Linker); ok {
				payload.Link(request.GeneratePageLinks(c))
				return json.Marshal(payload)
			}

			return json.Marshal(p)
		}
	}

	return nil, errors.New("no handler found supporting the requested format")
}
