package jsonapi

import "github.com/google/jsonapi"

type ManyPayload struct {
	jsonapi.ManyPayload
}

type Linker interface {
	Link(links jsonapi.Links)
}

func (p *ManyPayload) Link(links jsonapi.Links) {
	p.Links = &links
}
