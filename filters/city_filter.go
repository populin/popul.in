package filters

import (
	"net/url"

	"github.com/go-playground/form"
	"gopkg.in/olivere/elastic.v5"
)

// CityFilter holds the filter behavior and the query form data
type CityFilter struct {
	City bool `form:"city"`
}

// NewCityFilter is the factory method for CityFilter
func NewCityFilter() *CityFilter {
	return &CityFilter{City: false}
}

// Filter adds a TermQuery filter to the BoolQuery from city boolean value in the url
func (f *CityFilter) Filter(values url.Values, query *elastic.BoolQuery) error {
	decoder := form.NewDecoder()
	err := decoder.Decode(f, values)

	if values.Get("city") != "" {
		query.Filter(elastic.NewTermQuery("properties.city", f.City))
	}

	return err
}
