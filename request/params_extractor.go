package request

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
)

// SearchParamsExtractor contains the extracted params and the request to expose appropriate getters
type SearchParamsExtractor struct {
	params  *searchParams
	context *gin.Context
}

type searchParams struct {
	Longitude    float64    `form:"longitude"`
	Latitude     float64    `form:"latitude"`
	Radius       uint64     `form:"radius"`
	Geometry     bool       `form:"geometry"`
	Properties   Properties `form:"properties"`
	From         uint       `form:"from"`
	Size         uint       `form:"size"`
	SearchString string     `form:"q"`
}

// Properties represents properties enabled filters
type Properties struct {
	IsCity             bool   `form:"isCity"`
	AdministrativeName string `form:"administrativeName"`
	CountryCode        string `form:"country"`
}

// GetLatAndLon returns latitude, longitude and an error in case both were not found
func (e *SearchParamsExtractor) GetLatAndLon() (float64, float64, error) {
	_, latExists := e.context.GetQuery("latitude")
	_, lonExists := e.context.GetQuery("longitude")

	if latExists && lonExists {
		return e.params.Latitude, e.params.Longitude, nil
	}

	return 0, 0, fmt.Errorf("no latitude and longitude found")
}

// GetRadius returns the submitted radius (default 1)
func (e *SearchParamsExtractor) GetRadius() uint64 {
	if e.params.Radius == 0 {
		return 1
	}

	return e.params.Radius
}

// GetIsCity return isCity if it has been passed in the query
func (e *SearchParamsExtractor) GetIsCity() (bool, error) {
	if _, exists := e.context.GetQuery("properties.isCity"); !exists {
		return false, fmt.Errorf("no city filter")
	}

	return e.params.Properties.IsCity, nil
}

// GetGeometry returns the geometry boolean param
func (e *SearchParamsExtractor) GetGeometry() bool {
	return e.params.Geometry
}

// GetFrom returns the first item to return
func (e *SearchParamsExtractor) GetFrom() uint {
	return e.params.From
}

// GetSize returns the number of items to return
func (e *SearchParamsExtractor) GetSize() uint {
	return e.params.Size
}

// GetSearchString returns the search string
func (e *SearchParamsExtractor) GetSearchString() (string, error) {
	if _, exists := e.context.GetQuery("q"); !exists {
		return "", fmt.Errorf("no search string")
	}

	return e.params.SearchString, nil
}

// GetCountry returns the country string
func (e *SearchParamsExtractor) GetCountry() (string, error) {
	if _, exists := e.context.GetQuery("properties.country"); !exists {
		return "", fmt.Errorf("no country")
	}

	return e.params.Properties.CountryCode, nil
}

// GetAdministrativeName returns the type string
func (e *SearchParamsExtractor) GetAdministrativeName() (string, error) {
	if _, exists := e.context.GetQuery("properties.administrativeName"); !exists {
		return "", fmt.Errorf("no administrative name")
	}

	return e.params.Properties.AdministrativeName, nil
}

// NewSearchParamsExtractor returns an extractor from the Gin Context
func NewSearchParamsExtractor(c *gin.Context) (SearchParamsExtractor, error) {
	decoder := form.NewDecoder()

	params := searchParams{From: 0, Size: 50}

	err := decoder.Decode(&params, c.Request.URL.Query())

	if err != nil {
		return SearchParamsExtractor{}, err
	}

	return SearchParamsExtractor{params: &params, context: c}, err
}
