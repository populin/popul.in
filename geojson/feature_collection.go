package geojson

import (
	"encoding/json"

	"github.com/google/jsonapi"
)

// A FeatureCollection correlates to a GeoJSON feature collection.
type FeatureCollection struct {
	Type        string                 `json:"type"`
	BoundingBox []float64              `json:"bbox,omitempty"`
	Features    []*Feature             `json:"features"`
	CRS         map[string]interface{} `json:"crs,omitempty"`
	Links       *jsonapi.Links         `json:"links,omitempty"`
}

func (fc *FeatureCollection) Link(links jsonapi.Links) {
	fc.Links = &links
}

// NewFeatureCollection creates and initializes a new feature collection.
func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]*Feature, 0),
	}
}

// AddFeature appends a feature to the collection.
func (fc *FeatureCollection) AddFeature(feature *Feature) *FeatureCollection {
	fc.Features = append(fc.Features, feature)
	return fc
}

// MarshalJSON converts the feature collection object into the proper JSON.
// It will handle the encoding of all the child features and geometries.
// Alternately one can call json.Marshal(fc) directly for the same result.
func (fc *FeatureCollection) MarshalJSON() ([]byte, error) {
	fc.Type = "FeatureCollection"
	if fc.Features == nil {
		fc.Features = make([]*Feature, 0) // GeoJSON requires the feature attribute to be at least []
	}
	return json.Marshal(*fc)
}
