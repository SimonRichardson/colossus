package strategies

import (
	"github.com/SimonRichardson/colossus/cribs"
)

type hash struct{}

// NewHash defines a predetermined pool node depending on the hash of a key
func NewHash() *hash {
	return &hash{}
}

func (r *hash) Select(key string, max int) int {
	return int(cribs.New(key) % uint32(max))
}
