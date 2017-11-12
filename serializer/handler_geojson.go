package serializer

import (
	"encoding/json"

	"fmt"

	"github.com/populin/popul.in/constants"
	"github.com/populin/popul.in/geojson"
)

// GeoJSONHandler is an struct implementing FormatHandler interface
type GeoJSONHandler struct{}

// Supports defines which media type is supported (from FormatHandler interface)
func (GeoJSONHandler) Supports(format string) bool {
	return format == constants.GeoJSON
}

// Handle marshal the data (from FormatHandler interface)
func (GeoJSONHandler) Handle(o interface{}) ([]byte, error) {
	if features, ok := o.([]*geojson.Feature); ok {
		c := geojson.NewFeatureCollection()
		for _, feature := range features {
			c.AddFeature(feature)
		}
		r, err := json.Marshal(c)
		return r, err
	}

	if feature, ok := o.(*geojson.Feature); ok {
		r, err := json.Marshal(feature)
		return r, err
	}

	return []byte{}, fmt.Errorf("cannot handle type %T", o)
}