package main

import (
	"os"

	"github.com/populin/popul.in/api"
	"github.com/populin/popul.in/elastic"
)

func main() {
	if os.Getenv("ENV") != "dev" {
		api.SetReleaseMode()
	}

	clt, _ := elastic.NewClient()
	defer clt.Stop()

	app := api.Setup(clt)
	app.Run(":80")
}
