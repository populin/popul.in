package jsonapi

import "github.com/google/jsonapi"

// ManyPayload embed the original jsonapi.ManyPayload to implement the Linker interface
type ManyPayload struct {
	jsonapi.ManyPayload
}

// Linker interface describe the method used to add
type Linker interface {
	Link(links jsonapi.Links)
}

// Link add links to the ManyPayload
func (p *ManyPayload) Link(links jsonapi.Links) {
	p.Links = &links
}
