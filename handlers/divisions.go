package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/populin/popul.in/constants"
	"github.com/populin/popul.in/filters"
	"github.com/populin/popul.in/request"
	"github.com/populin/popul.in/serializer"
	"github.com/populin/popul.in/storage"
	"gopkg.in/olivere/elastic.v5"
)

// ByID returns a Feature by its ID
func ByID(c *gin.Context) {

	ds := c.MustGet("divisions_storage").(*storage.DivisionsStorage)

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
	ds := c.MustGet("divisions_storage").(*storage.DivisionsStorage)
	p := c.MustGet("pagination").(*request.Pagination)

	showGeometry := false
	var sorting elastic.Sorter

	if c.GetHeader("Accept") == constants.GeoJSON {
		showGeometry = true
		sorting = elastic.NewFieldSort("properties.administrativeLevel")
	} else {
		sorting = elastic.NewScoreSort()
	}

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
		return
	}

	features, total, err := ds.GetSearchResults(
		query,
		sorting,
		int((p.Page-1)*p.Size),
		int(p.Size),
		showGeometry,
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
