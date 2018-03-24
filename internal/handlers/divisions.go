package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"github.com/populin/popul.in/internal/constants"
	"github.com/populin/popul.in/internal/filters"
	"github.com/populin/popul.in/internal/request"
	"github.com/populin/popul.in/internal/serializer"
	"github.com/populin/popul.in/internal/storage"
)

// ByID returns a Feature by its ID
func ByID(c *gin.Context) {

	ds := c.MustGet("divisions_storage").(*storage.DivisionStorage)

	id := c.Param("id")

	showGeometry := c.GetHeader("Accept") == constants.GeoJSON

	feature, err := ds.FindOneByID(id, showGeometry)

	if err != nil {
		b := NewErrorBuilder()
		b.AddError(http.StatusNotFound, fmt.Sprintf("Division %s not found", id))
		handler := Error(b.Errors...)
		handler(c)
		return
	}

	s := serializer.NewSerializer()

	r, err := s.Serialize(c, feature)

	if err != nil {
		b := NewErrorBuilder()
		b.AddError(http.StatusBadRequest, err.Error())
		handler := Error(b.Errors...)
		handler(c)
		return
	}

	c.Writer.Write(r)
}

// Search parses the request to search for features
func Search(c *gin.Context) {
	ds := c.MustGet("divisions_storage").(*storage.DivisionStorage)
	p := c.MustGet("pagination").(*request.Pagination)

	var sorting elastic.Sorter

	if c.GetHeader("Accept") == constants.GeoJSON {
		sorting = elastic.NewFieldSort("properties.administrativeLevel")
	} else {
		sorting = elastic.NewScoreSort()
	}

	query := generateQuery(c)

	if c.IsAborted() {
		return
	}

	features, total, err := ds.GetSearchResults(
		query,
		sorting,
		int((p.Page-1)*p.Size),
		int(p.Size),
		c.GetHeader("Accept") == constants.GeoJSON,
	)

	p.TotalItems = uint(total)

	if err != nil {
		b := NewErrorBuilder()
		b.AddError(http.StatusServiceUnavailable, err.Error())
		handler := Error(b.Errors...)
		handler(c)
		return
	}

	s := serializer.NewSerializer()

	r, err := s.Serialize(c, features)

	if err != nil {
		b := NewErrorBuilder()
		b.AddError(http.StatusBadRequest, err.Error())
		handler := Error(b.Errors...)
		handler(c)
		return
	}

	c.Writer.Write(r)
}

func generateQuery(c *gin.Context) elastic.Query {
	f := filters.NewAggregator(
		c.Request.URL.Query(),
		filters.NewFeatureCoordinatesFilter(),
		filters.NewCityFilter(),
		filters.NewSearchFilter(),
	)

	query := elastic.NewBoolQuery()

	errs := f.Filter(query)

	if len(errs) > 0 {
		b := NewErrorBuilder()
		for _, err := range errs {
			b.AddError(http.StatusBadRequest, err.Error())
		}
		handler := Error(b.Errors...)
		handler(c)
	}

	return query
}
