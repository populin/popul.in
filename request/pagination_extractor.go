package request

import (
	"net/url"

	"github.com/go-playground/form"
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
