package serializer

import (
	"encoding/json"

	"github.com/google/jsonapi"
	"github.com/populin/popul.in/constants"
)

// JSONAPIHandler is an struct implementing FormatHandler interface
type JSONAPIHandler struct{}

// Supports defines which media type is supported (from FormatHandler interface)
func (JSONAPIHandler) Supports(format string) bool {
	return format == constants.JSONAPI
}

// Handle marshal the data (from FormatHandler interface)
func (JSONAPIHandler) Handle(o interface{}) ([]byte, error) {
	p, err := jsonapi.Marshal(o)

	if err != nil {
		return []byte{}, err
	}

	return json.Marshal(p)
}
