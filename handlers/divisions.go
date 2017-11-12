package handlers

import (
	"fmt"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/populin/popul.in/constants"
	"github.com/populin/popul.in/request"
	"github.com/populin/popul.in/serializer"
	"github.com/populin/popul.in/storage"
)

// ByID returns a Feature by its ID
func ByID(c *gin.Context) {

	ds := c.MustGet("DivisionsStorage").(*storage.DivisionsStorage)

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

// Search parses the request to search for a Feature collection
func Search(c *gin.Context) {
	ds := c.MustGet("DivisionsStorage").(*storage.DivisionsStorage)

	params, err := request.NewSearchParamsExtractor(c)

	if err != nil {
		decodeErrors := err.(form.DecodeErrors)
		ValidationError(http.StatusBadRequest, decodeErrors)(c)
		return
	}

	showGeometry := c.GetHeader("Accept") == constants.GeoJSON

	collection, total, err := ds.GetSearchResults(params, showGeometry)

	if err != nil {
		Error(http.StatusServiceUnavailable, err)(c)
		return
	}

	c.Header("X-Total-Results", strconv.FormatInt(total, 10))

	s := serializer.NewSerializer()

	r, err := s.Serialize(c, collection)

	if err != nil {
		Error(http.StatusBadRequest, err)(c)
		c.Abort()
		return
	}

	c.Writer.Write(r)
}
