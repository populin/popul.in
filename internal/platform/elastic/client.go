package elastic

import (
	"fmt"

	"github.com/olivere/elastic"
)

// NewClient is the ElasticSearch Client Factory
func NewClient(url string, port string) (*elastic.Client, error) {
	return elastic.NewSimpleClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%s", url, port)),
	)
}
