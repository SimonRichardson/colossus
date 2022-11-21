package common

import (
	"sort"

	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
)

type IntValue interface {
	Len() int
	Value() (int, error)
}

type SimilarInt struct {
	values map[int]struct{}
}

func NewSimilarInt() *SimilarInt {
	return &SimilarInt{make(map[int]struct{}, 0)}
}

func (x *SimilarInt) Similar(value int) bool {
	x.values[value] = struct{}{}

	return len(x.values) == 1
}

func (x *SimilarInt) Len() int {
	return len(x.values)
}

func (x *SimilarInt) Value() (int, error) {
	num := len(x.values)
	if num == 1 {
		for k := range x.values {
			return k, nil
		}
	}
	return 0, typex.Errorf(errors.Source, errors.UnexpectedResults,
		"Mismatch of similar values found (expected: 1, actual: %d).", num)
}

type LargestInt struct {
	values map[int]struct{}
}

func NewLargestInt() *LargestInt {
	return &LargestInt{make(map[int]struct{}, 0)}
}

func (x *LargestInt) Add(value int) {
	x.values[value] = struct{}{}
}

func (x *LargestInt) Len() int {
	return len(x.values)
}

func (x *LargestInt) Value() (int, error) {
	num := len(x.values)
	if num >= 1 {
		values := make([]int, 0, num)
		for k := range x.values {
			values = append(values, k)
		}
		sort.Ints(values)
		return values[len(values)-1], nil
	}
	return 0, typex.Errorf(errors.Source, errors.UnexpectedResults,
		"Mismatch of similar values found (expected: 1, actual: %d).", num)
}
