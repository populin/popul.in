package api

import (
	"fmt"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/populin/popul.in/handlers"
	"github.com/populin/popul.in/storage"
)

func byID(c *gin.Context) {

	ds := c.MustGet("DivisionsStorage").(*storage.DivisionsStorage)

	id := c.Param("id")

	showGeometry := true

	geometry, geometryExists := c.GetQuery("geometry")
	if geometryExists {
		showGeometry, _ = strconv.ParseBool(geometry)
	}

	feature, err := ds.FindOneByID(id, showGeometry)

	if err != nil {
		handlers.Error(http.StatusNotFound, fmt.Errorf("Division %s not found", id))(c)
		return
	}

	c.JSON(http.StatusOK, feature)

	return
}
