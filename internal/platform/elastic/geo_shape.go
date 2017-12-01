package elastic

import "fmt"

// GeoShapeCircleQuery is a simple geo_shape query implementation
type GeoShapeCircleQuery struct {
	longitude float64
	latitude  float64
	radius    uint64
}

// NewGeoShapeQuery creates and initializes a new geo shape query.
func NewGeoShapeQuery(lon float64, lat float64) *GeoShapeCircleQuery {
	return &GeoShapeCircleQuery{longitude: lon, latitude: lat, radius: 1}
}

// SetRadius set the radius of the Circle
func (q *GeoShapeCircleQuery) SetRadius(radius uint64) {
	q.radius = radius
}

// Source implements the Source interface
func (q GeoShapeCircleQuery) Source() (interface{}, error) {
	source := make(map[string]interface{})
	geoshape := make(map[string]interface{})
	geometry := make(map[string]interface{})
	shape := make(map[string]interface{})
	coordinates := [2]float64{q.longitude, q.latitude}

	shape["type"] = "circle"
	shape["radius"] = fmt.Sprintf("%dm", q.radius)
	shape["coordinates"] = coordinates
	geometry["shape"] = shape
	geoshape["geometry"] = geometry
	source["geo_shape"] = geoshape

	return source, nil
}
