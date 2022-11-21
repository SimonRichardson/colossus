package cribs

import (
	"bytes"
	"encoding/binary"
)

// Murmur3 implements the Murmur3 string hashing function. It can be passed to
// NewCluster.
//
// https://en.wikipedia.org/wiki/MurmurHash
const (
	c1 uint32 = 0xcc9e2d51
	c2 uint32 = 0x1b873593
	r1 uint32 = 15
	r2 uint32 = 13
	m  uint32 = 5
	n  uint32 = 0xe6546b64
)

// New returns a hash
func New(s string) uint32 {
	var (
		key    = []byte(s)
		length = len(key)
	)
	if length == 0 {
		return 0
	}

	var hash, k uint32
	buf := bytes.NewBufferString(s)

	nblocks := length / 4
	for i := 0; i < nblocks; i++ {
		binary.Read(buf, binary.LittleEndian, &k)

		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2

		hash ^= k
		hash = (hash << r2) | (hash >> (32 - r2))
		hash = (hash * m) + n
	}

	k = 0
	tailIndex := nblocks * 4
	switch length & 3 {
	case 3:
		k ^= uint32(key[tailIndex+2]) << 16
		fallthrough
	case 2:
		k ^= uint32(key[tailIndex+1]) << 8
		fallthrough
	case 1:
		k ^= uint32(key[tailIndex])

		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2

		hash ^= k
	}

	hash ^= uint32(length)
	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16

	return hash
}
