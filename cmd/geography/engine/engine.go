package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"github.com/populin/popul.in/cmd/geography/handlers"
	"github.com/populin/popul.in/internal/constants"
	generichandlers "github.com/populin/popul.in/internal/handlers"
	"github.com/populin/popul.in/internal/middlewares"
	"github.com/populin/popul.in/internal/serializer"
	"github.com/populin/popul.in/internal/storage"
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
		c.Set("divisions_storage", storage.New(ESClient))
		c.Set("serializer", serializer.NewSerializer(serializer.GeoJSONHandler{}, serializer.JSONAPIHandler{}))
		c.Next()
	})

	// error handlers
	router.NoRoute(generichandlers.NotFound())
	router.NoMethod(generichandlers.MethodNotAllowed())
	router.HandleMethodNotAllowed = true

	// global middlewares
	router.Use(gin.Recovery())
	router.Use(middlewares.CORS())

	divisions := router.Group("/divisions", middlewares.Negotiate(constants.GeoJSON, constants.JSONAPI))
	{
		divisions.GET("/:id", handlers.ByID)
		divisions.GET("", middlewares.Pagination(), handlers.Search)
	}

	return router
}
