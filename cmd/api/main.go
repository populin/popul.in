package main

import (
	"os"

	"github.com/populin/popul.in/cmd/api/engine"
	"github.com/populin/popul.in/internal/platform/elastic"
)

func main() {
	if os.Getenv("API_ENV") != "dev" {
		engine.SetReleaseMode()
	}

	clt, _ := elastic.NewClient(os.Getenv("ELASTIC_URL"), os.Getenv("ELASTIC_PORT"))
	defer clt.Stop()

	app := engine.Setup(clt)
	app.Run(":" + os.Getenv("API_PORT"))
}
