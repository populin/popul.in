package contexts

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/paulmach/go.geojson"
	"github.com/xeipuuv/gojsonschema"
)

// APIContext do the setup of the suite
func APIContext(s *godog.Suite, bs ...beforeScenario) {
	apiFeature := &APIFeature{}

	s.Step(`^I send a "(GET|POST|PUT|DELETE)" request to "([^"]*)" accepting "([^"]*)"$`, apiFeature.iSendARequestToAccepting)
	s.Step(`^the response code should be (\d+)$`, apiFeature.theResponseCodeShouldBe)
	s.Step(`^the response header "([^"]*)" should be "([^"]*)"$`, apiFeature.theResponseHeaderShouldBe)
	s.Step(`^the response header "([^"]*)" should exist$`, apiFeature.theResponseHeaderShouldExist)
	s.Step(`^the JSON response should be a valid GeoJson Feature$`, apiFeature.theJSONResponseShouldBeAValidGeoJSONFeature)
	s.Step(`^the JSON response should be a valid GeoJson Feature Collection$`, apiFeature.theJSONResponseShouldBeAValidGeoJSONFeatureCollection)
	s.Step(`^the GeoJSON property "([^"]*)" should be equal to "([^"]*)"$`, apiFeature.theGeoJSONPropertyShouldBeEqualTo)
	s.Step(`^the error message should be "([^"]*)"$`, apiFeature.theErrorMessageShouldBe)
	s.Step(`^the JSON should be valid according to this schema:$`, apiFeature.theJSONShouldBeValidAccordingToThisSchema)

	for _, b := range bs {
		s.BeforeScenario(b(apiFeature))
	}
}

// APIFeature holds a recorder and the engine to perform the tests
type APIFeature struct {
	Resp                     *httptest.ResponseRecorder
	App                      *gin.Engine
	geojsonFeature           *geojson.Feature
	geojsonFeatureCollection *geojson.FeatureCollection
}

type beforeScenario func(*APIFeature) func(interface{})

func (a *APIFeature) iSendARequestToAccepting(method string, endpoint string, accept string) (err error) {
	req, _ := http.NewRequest(method, endpoint, nil)

	req.Header.Set("Accept", accept)

	a.App.ServeHTTP(a.Resp, req)

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

func (a *APIFeature) theResponseCodeShouldBe(code int) error {
	if code != a.Resp.Code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.Resp.Code)
	}

	return nil
}

func (a *APIFeature) theResponseHeaderShouldBe(header string, value string) error {
	if headerValue := a.Resp.Result().Header.Get(header); headerValue != value {
		return fmt.Errorf("expected response header %s to be: %s, but actual is: %s", header, value, headerValue)
	}

	return nil
}

func (a *APIFeature) theResponseHeaderShouldExist(header string) error {
	if headerValue := a.Resp.Result().Header.Get(header); headerValue == "" {
		return fmt.Errorf("response header %s was not found", header)
	}

	return nil
}

func (a *APIFeature) theJSONResponseShouldBeAValidGeoJSONFeature() error {
	var err error
	a.geojsonFeature, err = geojson.UnmarshalFeature(a.Resp.Body.Bytes())

	if err != nil {
		return fmt.Errorf("the response could not be unserialized to a GeoJSON Feature")
	}

	return nil
}

func (a *APIFeature) theJSONResponseShouldBeAValidGeoJSONFeatureCollection() error {
	var err error
	a.geojsonFeatureCollection, err = geojson.UnmarshalFeatureCollection(a.Resp.Body.Bytes())

	if err != nil {
		return fmt.Errorf("the response could not be unserialized to a GeoJSON Feature Collection")
	}
	return nil
}

func (a *APIFeature) theGeoJSONPropertyShouldBeEqualTo(prop string, value string) error {
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

func (a *APIFeature) theErrorMessageShouldBe(message string) error {
	errorsPayload := jsonapi.ErrorsPayload{}

	err := json.Unmarshal(a.Resp.Body.Bytes(), &errorsPayload)

	if err != nil {
		return err
	}

	var messages []string

	for _, err := range errorsPayload.Errors {
		if err.Detail == message {
			return nil
		}
		messages = append(messages, err.Detail)
	}

	return fmt.Errorf("message error \"%s\" was not found (found: %s)", message, strings.Join(messages, ", "))
}

func (a *APIFeature) theJSONShouldBeValidAccordingToThisSchema(body *gherkin.DocString) error {
	schemaLoader := gojsonschema.NewStringLoader(body.Content)
	documentLoader := gojsonschema.NewStringLoader(a.Resp.Body.String())

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
