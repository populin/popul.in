package serializer

import (
	"fmt"

	"github.com/gin-gonic/gin"
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
func (GeoJSONHandler) Handle(c *gin.Context, o interface{}) (interface{}, error) {
	if features, ok := o.([]*geojson.Feature); ok {
		c := geojson.NewFeatureCollection()
		for _, feature := range features {
			c.AddFeature(feature)
		}
		return c, nil
	}

	if feature, ok := o.(*geojson.Feature); ok {
		return feature, nil
	}

	return nil, fmt.Errorf("cannot handle type %T", o)
}
