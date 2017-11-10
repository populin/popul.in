package constants

const (
	// Index is the ElasticSearch Index
	Index = "geography"
	// Type is the ElasticSearch Type
	Type = "division"
	// Mapping is the ElasticSearch Mapping for divisions in geography
	Mapping = `
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
)
