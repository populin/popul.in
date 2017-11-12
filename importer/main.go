package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/paulmach/go.geojson"
	"github.com/populin/popul.in/constants"
	es "github.com/populin/popul.in/elastic"
	"github.com/populin/popul.in/slugger"
	"github.com/urfave/cli"
	"gopkg.in/olivere/elastic.v5"
)

// Mapping is the ElasticSearch Mapping for divisions in geography
const Mapping = `
	{
		"settings":{
			"number_of_shards": 8,
			"number_of_replicas": 0
		},
		"mappings":{
			"division":{
				"properties":{
					"geometry": {
						"type": "geo_shape"
					},
					"properties": {
						"properties": {
							"name": {
							  "type": "string"
							},
							"administrativeLevel": {
							  "type": "integer"
							},
							"administrativeName": {
							  "type": "text",
							  "fielddata": true
							},
							"code": {
							  "type": "string"
							},
							"country": {
							  "type": "string"
							},
							"isCity": {
							  "type": "boolean"
							},
							"isCountry": {
							  "type": "boolean"
							}
						}
					}
				}
			}
		}
	}`

func main() {
	app := cli.NewApp()
	app.Name = "importer"
	app.Usage = "geojson files parser and importer"
	app.Action = importFolder

	app.Run(os.Args)
}

func importFolder(c *cli.Context) error {

	folder := c.Args().Get(0)

	files, err := filepath.Glob(fmt.Sprintf("%s/*.json", folder))
	if err != nil {
		return err
	}

	processor, err := getProcessor()

	if err != nil {
		return err
	}

	for _, file := range files {

		f, _ := ioutil.ReadFile(file)
		collection, err := geojson.UnmarshalFeatureCollection(f)

		if err != nil {
			return err
		}

		fmt.Printf("Parsing file %s => %d features\n", file, len(collection.Features))

		for _, feature := range collection.Features {
			req := elastic.NewBulkIndexRequest().
				Index(constants.ESIndexGeography).
				Type(constants.ESTypeGeography).
				Id(sluggify(feature)).
				Doc(feature)

			processor.Add(req)
		}
	}

	errClose := processor.Close()

	return errClose
}

func getProcessor() (*elastic.BulkProcessor, error) {
	ctx := context.Background()

	clt, err := es.NewClient()

	if err != nil {
		return nil, err
	}

	exists, err := clt.IndexExists(constants.ESIndexGeography).Do(ctx)

	if err != nil {
		panic(err)
	}

	if exists {
		deleteIndex, err := clt.DeleteIndex(constants.ESIndexGeography).Do(ctx)
		if err != nil {
			return nil, err
		}
		if !deleteIndex.Acknowledged {
			// Not acknowledged
		}
	}

	createIndex, err := clt.CreateIndex(constants.ESIndexGeography).
		BodyString(Mapping).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	if !createIndex.Acknowledged {
		// Not acknowledged
	}

	return clt.BulkProcessor().
		Name("GeographyWorker").
		BulkSize(-1).
		BulkActions(1000).
		Do(context.Background())
}

func sluggify(feature *geojson.Feature) string {

	slug, err := slugger.Sluggify(
		[]string{
			feature.PropertyMustString("country"),
			feature.PropertyMustString("administrativeName"),
			feature.PropertyMustString("name"),
		},
		feature.Geometry,
	)

	if err != nil {
		panic(err)
	}

	return slug
}
