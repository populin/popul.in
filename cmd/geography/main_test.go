package main

import (
	"os"
	"testing"
	"time"

	"net/http/httptest"

	"github.com/DATA-DOG/godog"
	"github.com/populin/popul.in/cmd/geography/engine"
	"github.com/populin/popul.in/internal/contexts"
	"github.com/populin/popul.in/internal/platform/elastic"
)

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("geography", func(s *godog.Suite) {
		beforeScenario := func(a *contexts.APIFeature) func(interface{}) {
			return func(interface{}) {
				engine.SetTestMode()

				a.Resp = httptest.NewRecorder()

				clt, _ := elastic.NewClient(os.Getenv("POPULIN_GEOGRAPHY_ELASTIC_URL"), os.Getenv("POPULIN_GEOGRAPHY_ELASTIC_PORT"))

				a.App = engine.Setup(clt)
			}
		}

		contexts.APIContext(s, beforeScenario)
	}, godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(),
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
