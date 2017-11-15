package constants

import "github.com/google/jsonapi"

const (
	// GeoJSON is the identifier for the GeoJSON media type
	GeoJSON = "application/geo+json"
	// JSONAPI is the identifier for the JSON API media type
	JSONAPI = jsonapi.MediaType
	// ESIndexGeography is the ElasticSearch Index
	ESIndexGeography = "geography"
	// ESTypeGeography is the ElasticSearch Type
	ESTypeGeography = "division"
)
