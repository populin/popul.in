package main

import (
	"os"

	"github.com/populin/popul.in/cmd/politics/engine"
)

func main() {
	if os.Getenv("POPULIN_POLITICS_API_ENV") != "dev" {
		engine.SetReleaseMode()
	}

	app := engine.Setup()
	app.Run(":" + os.Getenv("POPULIN_POLITICS_API_PORT"))
}
