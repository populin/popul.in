package filters

import (
	"net/url"

	"github.com/go-playground/form"
	"gopkg.in/olivere/elastic.v5"
)

// ElasticFilter defines the Filter method which builds the ElasticSearch BoolQuery depending on url values
type ElasticFilter interface {
	// Filter add new conditions to the BoolQuery from url
	Filter(values url.Values, query *elastic.BoolQuery) error
}

// Aggregator is the registry of ElasticFilters
type Aggregator struct {
	values  url.Values
	filters []ElasticFilter
}

// NewAggregator is the factory method of Aggregator
func NewAggregator(values url.Values, filters ...ElasticFilter) Aggregator {
	return Aggregator{values: values, filters: filters}
}

// Filter calls all registered filters Filter methods
func (a Aggregator) Filter(query *elastic.BoolQuery) map[string]error {
	errs := make(map[string]error)

	for _, filter := range a.filters {
		err := filter.Filter(a.values, query)
		if err != nil {
			if decodeErrors, ok := err.(form.DecodeErrors); ok {
				return decodeErrors
			}
			errs["Bad Request"] = err
		}
	}

	return errs
}
