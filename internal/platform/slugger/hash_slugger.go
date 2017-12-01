package slugger

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
	"github.com/mitchellh/hashstructure"
)

// Sluggify creates a slug with a slice of strings and an interface-based unicity
func Sluggify(parts []string, hashBase interface{}) (string, error) {

	hash, err := hashstructure.Hash(hashBase, nil)

	if err != nil {
		return "", err
	}

	parts = append(parts, fmt.Sprintf("%d", hash))

	str := strings.Join(parts, " ")

	return slug.Make(str), nil
}
