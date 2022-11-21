package common

import (
	"fmt"
	"strings"
)

// Normalise defines a way to normalise a string, but is commonly users for
// normalising environmental variables.
func Normalise(strategy string) string {
	return StripWhitespace(strings.ToLower(strategy))
}

// StripWhitespace defines a way to remove newlines, tabs and spaces from a
// string, this isn't bullet proof and not expected to work in every location.
// But for reading in environmental variables this should suffice.
func StripWhitespace(src string) string {
	var dst []rune
	for _, c := range src {
		switch c {
		case ' ', '\t', '\r', '\n':
			continue
		}
		dst = append(dst, c)
	}
	return string(dst)
}

func SumIntegers(values []int) int {
	res := 0
	for _, v := range values {
		res += v
	}
	return res
}

func AvgIntegers(values []int) int {
	num := len(values)
	if num < 1 {
		return 0
	}
	return SumIntegers(values) / num
}

func SumErrors(values []error) error {
	var (
		filtered = FilterErrors(values)
		result   = make([]string, len(filtered))
	)
	for k, v := range filtered {
		result[k] = v.Error()
	}
	return fmt.Errorf(strings.Join(result, "; "))
}

func FilterErrors(values []error) []error {
	result := make([]error, 0, len(values))
	for _, v := range values {
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}

func LargestInteger(values []int) int {
	res := 0
	for _, v := range values {
		if v > res {
			res = v
		}
	}
	return res
}
