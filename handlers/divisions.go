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
		Error(http.StatusNotFound, fmt.Errorf("Division %s not found", id))(c)
		c.Abort()
		return
	}

	s := serializer.NewSerializer()

	r, err := s.Serialize(c, feature)

	if err != nil {
		Error(http.StatusBadRequest, err)(c)
		c.Abort()
		return
	}

	c.Writer.Write(r)
}

// Search parses the request to search for features
func Search(c *gin.Context) {
	ds := c.MustGet("divisions_storage").(*storage.DivisionsStorage)
	p := c.MustGet("pagination").(*request.Pagination)

	showGeometry := c.GetHeader("Accept") == constants.GeoJSON

	f := filters.NewAggregator(
		c.Request.URL.Query(),
		filters.NewFeatureCoordinatesFilter(),
		filters.NewCityFilter(),
		filters.NewSearchFilter(),
	)

	query := elastic.NewBoolQuery()

	errs := f.Filter(query)

	if len(errs) > 0 {
		ValidationError(http.StatusBadRequest, errs)(c)
		c.Abort()
		return
	}

	features, total, err := ds.GetSearchResults(
		query,
		int((p.Page-1)*p.Size),
		int(p.Size),
		showGeometry,
	)

	p.TotalItems = uint(total)

	if err != nil {
		Error(http.StatusServiceUnavailable, err)(c)
		c.Abort()
		return
	}

	s := serializer.NewSerializer()

	r, err := s.Serialize(c, features)

	if err != nil {
		Error(http.StatusBadRequest, err)(c)
		c.Abort()
		return
	}

	c.Writer.Write(r)
}
