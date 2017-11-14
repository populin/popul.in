package filters

import (
	"net/url"

	"github.com/go-playground/form"
	"github.com/pkg/errors"
	es "github.com/populin/popul.in/elastic"
	"gopkg.in/olivere/elastic.v5"
)

// FeatureCoordinatesFilter holds the filter behavior and the query form data
type FeatureCoordinatesFilter struct {
	Longitude float64 `form:"longitude"`
	Latitude  float64 `form:"latitude"`
	Radius    uint64  `form:"radius"`
}

// NewFeatureCoordinatesFilter is the factory method for FeatureCoordinatesFilter
func NewFeatureCoordinatesFilter() *FeatureCoordinatesFilter {
	return &FeatureCoordinatesFilter{Radius: 10}
}

// Filter adds a GeoShapeQuery to the BoolQuery from coordinates in the url
func (f *FeatureCoordinatesFilter) Filter(values url.Values, query *elastic.BoolQuery) error {
	decoder := form.NewDecoder()
	err := decoder.Decode(f, values)

	if values.Get("latitude") != "" && values.Get("longitude") != "" {
		q := es.NewGeoShapeQuery(f.Longitude, f.Latitude)
		q.SetRadius(f.Radius)

		query.Must(q)
	}

	if values.Get("latitude") != "" && values.Get("longitude") == "" {
		return errors.New("missing longitude")
	}

	if values.Get("longitude") != "" && values.Get("latitude") == "" {
		return errors.New("missing latitude")
	}

	return err
}
