package main

import (
	"os"
	"testing"
	"time"

	"net/http/httptest"

	"github.com/DATA-DOG/godog"
	"github.com/populin/popul.in/cmd/politics/engine"
	"github.com/populin/popul.in/internal/contexts"
)

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("politics", func(s *godog.Suite) {
		beforeScenario := func(a *contexts.APIFeature) func(interface{}) {
			return func(interface{}) {
				engine.SetTestMode()

				a.Resp = httptest.NewRecorder()

				a.App = engine.Setup()
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
