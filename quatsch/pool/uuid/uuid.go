package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync/atomic"

	"github.com/SimonRichardson/colossus/quatsch/pool"
)

type UUID [24]byte

func (u UUID) Bytes() [24]byte {
	return [24]byte(u)
}

func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

type generator struct {
	seed    [24]byte
	counter uint64
}

func NewGenerator() (*generator, error) {
	var gen generator
	if _, err := rand.Read(gen.seed[:]); err != nil {
		return nil, pool.ErrInvalidSeed
	}
	return &gen, nil
}

func (g *generator) Next() UUID {
	var (
		counterBytes [8]byte
		x            = atomic.AddUint64(&g.counter, 1)
	)

	binary.LittleEndian.PutUint64(counterBytes[:], x)

	uuid := g.seed
	for i, b := range counterBytes {
		uuid[i] ^= b
	}
	return UUID(uuid)
}
