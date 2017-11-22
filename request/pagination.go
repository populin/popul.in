package request

import (
	"net/url"

	"math"
	strconv "strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/form"
	"github.com/google/jsonapi"
)

// Pagination holds values extracted from the request
type Pagination struct {
	Page       uint `form:"page[number]"`
	Size       uint `form:"page[size]"`
	TotalItems uint
}

// ExtractPagination gets the pagination params from the Gin Context
func ExtractPagination(v url.Values) (*Pagination, error) {
	decoder := form.NewDecoder()

	p := Pagination{Page: 1, Size: 50}

	err := decoder.Decode(&p, v)

	if err != nil {
		return nil, err
	}

	if p.Page < 1 {
		p.Page = 1
	}

	if p.Size > 500 {
		p.Size = 500
	}

	if p.Size == 0 {
		p.Size = 1
	}

	return &p, err
}

// GeneratePageLinks returns un valid jsonapi.Links struct with page links
func GeneratePageLinks(c *gin.Context) jsonapi.Links {
	pi, found := c.Get("pagination")

	links := jsonapi.Links{}

	if !found {
		return links
	}

	p := pi.(*Pagination)

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

	return links
}

func generatePageURI(u *url.URL, number uint, size uint) string {
	q := u.Query()

	q.Set("page[number]", strconv.Itoa(int(number)))
	q.Set("page[size]", strconv.Itoa(int(size)))

	u.RawQuery = q.Encode()

	return u.String()
}
