package serializer

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/populin/popul.in/internal/constants"
	ja "github.com/populin/popul.in/internal/platform/jsonapi"
)

// JSONAPIHandler is an struct implementing FormatHandler interface
type JSONAPIHandler struct{}

// Supports defines which media type is supported (from FormatHandler interface)
func (JSONAPIHandler) Supports(format string) bool {
	return format == constants.JSONAPI
}

// Handle marshal the data (from FormatHandler interface)
func (JSONAPIHandler) Handle(c *gin.Context, o interface{}) (interface{}, error) {
	p, err := jsonapi.Marshal(o)

	if payload, ok := p.(*jsonapi.ManyPayload); ok {
		newPayload := ja.ManyPayload{ManyPayload: *payload}

		return &newPayload, nil
	}

	return p, err
}
