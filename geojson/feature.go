package geojson

import (
	"encoding/json"

	"github.com/paulmach/go.geojson"
)

// A Feature corresponds to GeoJSON feature object
type Feature struct {
	ID          string                 `jsonapi:"primary,features" json:"id"`
	Type        string                 `jsonapi:"attr,type" json:"type"`
	BoundingBox []float64              `json:"bbox,omitempty"`
	Geometry    *geojson.Geometry      `json:"geometry,omitempty"`
	Properties  map[string]interface{} `jsonapi:"attr,properties" json:"properties"`
	CRS         map[string]interface{} `json:"crs,omitempty"` // Coordinate Reference System Objects are not currently supported
}

// MarshalJSON converts the feature object into the proper JSON.
// It will handle the encoding of all the child geometries.
// Alternately one can call json.Marshal(f) directly for the same result.
func (f *Feature) MarshalJSON() ([]byte, error) {
	f.Type = "Feature"
	if len(f.Properties) == 0 {
		f.Properties = nil
	}

	return json.Marshal(*f)
}

// UnmarshalFeature decodes the data into a GeoJSON feature.
// Alternately one can call json.Unmarshal(f) directly for the same result.
func UnmarshalFeature(data []byte) (*Feature, error) {
	f := &Feature{}
	err := json.Unmarshal(data, f)
	if err != nil {
		return nil, err
	}

	return f, nil
}
