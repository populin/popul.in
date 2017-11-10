package elastic

import (
	"fmt"
	"os"

	"gopkg.in/olivere/elastic.v5"
)

// NewClient is the ElasticSearch Client Factory
func NewClient() (*elastic.Client, error) {
	url := os.Getenv("ELASTIC_URL")
	port := os.Getenv("ELASTIC_PORT")

	return elastic.NewSimpleClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%s", url, port)),
	)
}
