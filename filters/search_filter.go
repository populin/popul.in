package filters

import (
	"net/url"

	"github.com/go-playground/form"
	"gopkg.in/olivere/elastic.v5"
)

// SearchFilter holds the filter behavior and the query form data
type SearchFilter struct {
	Search string `form:"q"`
}

// NewSearchFilter is the factory method for SearchFilter
func NewSearchFilter() *SearchFilter {
	return &SearchFilter{}
}

// Filter adds a MultiMatchQuery to the BoolQuery from search string in the url
func (f *SearchFilter) Filter(values url.Values, query *elastic.BoolQuery) error {
	decoder := form.NewDecoder()
	err := decoder.Decode(f, values)

	if values.Get("q") != "" {
		query.Must(elastic.NewMultiMatchQuery(f.Search, "properties.name^3", "properties.code^2"))
	}

	return err
}
