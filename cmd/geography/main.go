package main

import (
	"os"

	"github.com/populin/popul.in/cmd/geography/engine"
	"github.com/populin/popul.in/internal/platform/elastic"
)

func main() {
	if os.Getenv("POPULIN_GEOGRAPHY_API_ENV") != "dev" {
		engine.SetReleaseMode()
	}

	clt, _ := elastic.NewClient(os.Getenv("POPULIN_GEOGRAPHY_ELASTIC_URL"), os.Getenv("POPULIN_GEOGRAPHY_ELASTIC_PORT"))
	defer clt.Stop()

	app := engine.Setup(clt)
	app.Run(":" + os.Getenv("POPULIN_GEOGRAPHY_API_PORT"))
}
