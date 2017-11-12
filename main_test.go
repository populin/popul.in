package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gin-gonic/gin"
	"github.com/paulmach/go.geojson"
	"github.com/pkg/errors"
	"github.com/populin/popul.in/api"
	"github.com/populin/popul.in/elastic"
	"github.com/populin/popul.in/handlers"
	"github.com/xeipuuv/gojsonschema"
)

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
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

type apiFeature struct {
	resp                     *httptest.ResponseRecorder
	app                      *gin.Engine
	geojsonFeature           *geojson.Feature
	geojsonFeatureCollection *geojson.FeatureCollection
}

func (a *apiFeature) resetResponse(interface{}) {
	api.SetTestMode()

	a.resp = httptest.NewRecorder()

	clt, _ := elastic.NewClient()

	a.app = api.Setup(clt)
}

func (a *apiFeature) iSendARequestToAccepting(method string, endpoint string, accept string) (err error) {
	req, _ := http.NewRequest(method, endpoint, nil)

	req.Header.Set("Accept", accept)

	a.app.ServeHTTP(a.resp, req)

	// handle panic
	defer func() {
		switch t := recover().(type) {
		case string:
			err = fmt.Errorf(t)
		case error:
			err = t
		}
	}()

	return
}

func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.Code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
	}

	return nil
}

func (a *apiFeature) theResponseHeaderShouldBe(header string, value string) error {
	if headerValue := a.resp.Result().Header.Get(header); headerValue != value {
		return fmt.Errorf("expected response header %s to be: %s, but actual is: %s", header, value, headerValue)
	}

	return nil
}

func (a *apiFeature) theResponseHeaderShouldExist(header string) error {
	if headerValue := a.resp.Result().Header.Get(header); headerValue == "" {
		return fmt.Errorf("response header %s was not found", header)
	}

	return nil
}

func (a *apiFeature) theJSONResponseShouldBeAValidGeoJSONFeature() error {
	var err error
	a.geojsonFeature, err = geojson.UnmarshalFeature(a.resp.Body.Bytes())

	if err != nil {
		return fmt.Errorf("the response could not be unserialized to a GeoJSON Feature")
	}

	return nil
}

func (a *apiFeature) theJSONResponseShouldBeAValidGeoJSONFeatureCollection() error {
	var err error
	a.geojsonFeatureCollection, err = geojson.UnmarshalFeatureCollection(a.resp.Body.Bytes())

	if err != nil {
		return fmt.Errorf("the response could not be unserialized to a GeoJSON Feature Collection")
	}
	return nil
}

func (a *apiFeature) theGeoJSONPropertyShouldBeEqualTo(prop string, value string) error {
	propValue, err := a.geojsonFeature.PropertyString(prop)

	if err != nil {
		for k, v := range a.geojsonFeature.Properties {
			fmt.Printf("DEBUG: property %s => (%v, %T)\n", k, v, v)
		}
		return fmt.Errorf("error getting the property %s as string (%s)", prop, err.Error())
	}

	if propValue != value {
		return fmt.Errorf("the value of the property %s is %s, not %s", prop, propValue, value)
	}

	return nil
}

func (a *apiFeature) theErrorMessageShouldBe(message string) error {
	errorMessage := handlers.ErrorResponse{}

	err := json.Unmarshal(a.resp.Body.Bytes(), &errorMessage)

	if err != nil {
		return err
	}

	if errorMessage.Message != message {
		return fmt.Errorf("expected error message to be %s, but actual is %s", message, errorMessage.Message)
	}

	return nil
}

func (a *apiFeature) theJSONShouldBeValidAccordingToThisSchema(body *gherkin.DocString) error {
	schemaLoader := gojsonschema.NewStringLoader(body.Content)
	documentLoader := gojsonschema.NewStringLoader(a.resp.Body.String())

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		err := errors.New("the document is not valid")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return err
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	apiFeature := &apiFeature{}

	s.BeforeScenario(apiFeature.resetResponse)

	s.Step(`^I send a "(GET|POST|PUT|DELETE)" request to "([^"]*)" accepting "([^"]*)"$`, apiFeature.iSendARequestToAccepting)
	s.Step(`^the response code should be (\d+)$`, apiFeature.theResponseCodeShouldBe)
	s.Step(`^the response header "([^"]*)" should be "([^"]*)"$`, apiFeature.theResponseHeaderShouldBe)
	s.Step(`^the response header "([^"]*)" should exist$`, apiFeature.theResponseHeaderShouldExist)
	s.Step(`^the JSON response should be a valid GeoJson Feature$`, apiFeature.theJSONResponseShouldBeAValidGeoJSONFeature)
	s.Step(`^the JSON response should be a valid GeoJson Feature Collection$`, apiFeature.theJSONResponseShouldBeAValidGeoJSONFeatureCollection)
	s.Step(`^the GeoJSON property "([^"]*)" should be equal to "([^"]*)"$`, apiFeature.theGeoJSONPropertyShouldBeEqualTo)
	s.Step(`^the error message should be "([^"]*)"$`, apiFeature.theErrorMessageShouldBe)
	s.Step(`^the JSON should be valid according to this schema:$`, apiFeature.theJSONShouldBeValidAccordingToThisSchema)
}
