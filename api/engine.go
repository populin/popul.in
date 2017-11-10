package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/populin/popul.in/handlers"
	"github.com/populin/popul.in/middlewares"
	"github.com/populin/popul.in/storage"
	"gopkg.in/olivere/elastic.v5"
)

var modeName = gin.DebugMode

// SetReleaseMode switches Gin to PROD environment
func SetReleaseMode() {
	modeName = gin.ReleaseMode
}

// SetTestMode switches Gin to TEST environment
func SetTestMode() {
	modeName = gin.TestMode
}

// Setup configure the Gin Engine
func Setup(ESClient *elastic.Client) *gin.Engine {

	gin.SetMode(modeName)

	router := gin.New()

	if modeName == gin.DebugMode {
		router.Use(gin.Logger())
	}

	router.Use(func(c *gin.Context) {
		c.Set("DivisionsStorage", storage.New(ESClient))

		c.Next()
	})

	// error handlers
	router.NoRoute(handlers.Error(http.StatusNotFound, errors.New("Not found")))
	router.NoMethod(handlers.Error(http.StatusMethodNotAllowed, errors.New("Method not allowed")))
	router.HandleMethodNotAllowed = true

	// middlewares
	router.Use(middlewares.CORS())
	router.Use(middlewares.GeoJSON())

	divisions := router.Group("/divisions")
	{
		divisions.GET("/:id", byID)
		divisions.GET("", search)
	}

	return router
}