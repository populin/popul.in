package serializer

import (
	"encoding/json"

	"math"

	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/populin/popul.in/constants"
	"github.com/populin/popul.in/request"
)

// JSONAPIHandler is an struct implementing FormatHandler interface
type JSONAPIHandler struct{}

// Supports defines which media type is supported (from FormatHandler interface)
func (JSONAPIHandler) Supports(format string) bool {
	return format == constants.JSONAPI
}

// Handle marshal the data (from FormatHandler interface)
func (JSONAPIHandler) Handle(c *gin.Context, o interface{}) ([]byte, error) {
	p, err := jsonapi.Marshal(o)

	if payload, ok := p.(*jsonapi.ManyPayload); ok {
		payload.Links = generatePageLinks(c)
	}

	if err != nil {
		return []byte{}, err
	}

	return json.Marshal(p)
}

func generatePageLinks(c *gin.Context) *jsonapi.Links {
	p := c.MustGet("pagination").(*request.Pagination)
	links := jsonapi.Links{}

	last := uint(math.Ceil(float64(p.TotalItems) / float64(p.Size)))

	prev := p.Page - 1
	if prev < 1 {
		prev = 1
	}

	next := p.Page + 1
	if next > last {
		next = last
	}

	u := c.Request.URL

	links["self"] = generatePageURI(u, p.Page, p.Size)
	links["first"] = generatePageURI(u, 1, p.Size)
	links["last"] = generatePageURI(u, last, p.Size)
	links["prev"] = generatePageURI(u, prev, p.Size)
	links["next"] = generatePageURI(u, next, p.Size)

	return &links
}

func generatePageURI(u *url.URL, number uint, size uint) string {
	q := u.Query()

	q.Set("page[number]", strconv.Itoa(int(number)))
	q.Set("page[size]", strconv.Itoa(int(size)))

	u.RawQuery = q.Encode()

	return u.String()
}
