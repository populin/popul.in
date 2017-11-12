package storage

import (
	"context"
	"encoding/json"

	"github.com/populin/popul.in/constants"
	es "github.com/populin/popul.in/elastic"
	"github.com/populin/popul.in/geojson"
	"github.com/populin/popul.in/request"
	"gopkg.in/olivere/elastic.v5"
)

// New is the DivisionsStorage factory
func New(client *elastic.Client) *DivisionsStorage {
	return &DivisionsStorage{client}
}

// DivisionsStorage embed the ES client and is the gateway for stored division
type DivisionsStorage struct {
	client *elastic.Client
}

// FindOneByID gets a division by its ID
func (storage *DivisionsStorage) FindOneByID(id string, showGeometry bool) (*geojson.Feature, error) {
	var fsc *elastic.FetchSourceContext

	if !showGeometry {
		fsc = elastic.NewFetchSourceContext(true).Exclude("geometry")
	}

	division, err := storage.client.Get().
		Index(constants.ESIndexGeography).
		Type(constants.ESTypeGeography).
		FetchSourceContext(fsc).
		Id(id).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	feature, err := unmarshalFeature(*division.Source, division.Id)

	if err != nil {
		return nil, err
	}

	return feature, nil
}

// GetGeoShapeQuery returns a geo_shape Query
func (storage *DivisionsStorage) GetGeoShapeQuery(lat float64, lon float64, radius uint64) elastic.Query {
	query := es.NewGeoShapeQuery(lon, lat)
	query.SetRadius(radius)

	return query
}

// GetSearchResults returns a FeatureCollection from a BoolQuery
func (storage *DivisionsStorage) GetSearchResults(params request.SearchParamsExtractor, showGeometry bool) ([]*geojson.Feature, int64, error) {
	query := elastic.NewBoolQuery()

	if lat, lon, err := params.GetLatAndLon(); err == nil {
		radius := params.GetRadius()
		query.Must(storage.GetGeoShapeQuery(lat, lon, radius))
	}

	if isCity, err := params.GetIsCity(); err == nil {
		query.Filter(elastic.NewTermQuery("properties.isCity", isCity))
	}

	if searchString, err := params.GetSearchString(); err == nil {
		query.Must(elastic.NewMultiMatchQuery(searchString, "properties.name^3", "properties.code^2"))
	}

	if administrativeName, err := params.GetAdministrativeName(); err == nil {
		query.Must(elastic.NewMatchQuery("properties.administrativeName", administrativeName))
	}

	if country, err := params.GetCountry(); err == nil {
		query.Must(elastic.NewMatchPhraseQuery("properties.country", country))
	}

	var fsc *elastic.FetchSourceContext
	if !showGeometry {
		fsc = elastic.NewFetchSourceContext(true).Exclude("geometry")
	}

	results, err := storage.client.Search().
		Index(constants.ESIndexGeography).
		Type(constants.ESTypeGeography).
		From(int(params.GetFrom())).Size(int(params.GetSize())).
		Query(query).
		FetchSourceContext(fsc).
		Do(context.Background())

	if err != nil {
		return nil, 0, err
	}

	collection, err := unmarshalFeatureCollection(results.Hits.Hits)

	return collection, results.TotalHits(), err
}

func unmarshalFeatureCollection(hits []*elastic.SearchHit) ([]*geojson.Feature, error) {
	var features []*geojson.Feature
	for _, division := range hits {
		feature, err := unmarshalFeature(*division.Source, division.Id)

		if err != nil {
			return features, err
		}

		features = append(features, feature)
	}

	return features, nil
}

func unmarshalFeature(message json.RawMessage, id string) (*geojson.Feature, error) {
	feature, err := geojson.UnmarshalFeature(message)

	feature.ID = id

	if err != nil {
		return nil, err
	}

	return feature, nil
}
