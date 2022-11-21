package common

import "strings"

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
