package api

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/populin/popul.in/handlers"
	"github.com/populin/popul.in/request"
	"github.com/populin/popul.in/storage"
)

func search(c *gin.Context) {
	ds := c.MustGet("DivisionsStorage").(*storage.DivisionsStorage)

	params, err := request.NewSearchParamsExtractor(c)

	if err != nil {
		decodeErrors := err.(form.DecodeErrors)
		handlers.ValidationError(http.StatusBadRequest, decodeErrors)(c)
		return
	}

	collection, total, err := ds.GetSearchResults(params)

	if err != nil {
		handlers.Error(http.StatusServiceUnavailable, err)(c)
		return
	}

	c.Header("X-Total-Results", strconv.FormatInt(total, 10))

	c.JSON(http.StatusOK, collection)
	c.Abort()

	return
}
