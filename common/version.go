package common

import (
	"strconv"
	"strings"

	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
)

const (
	numbers string = "0123456789"
)

// ParseSemver parses version string and returns a validated Version or error
func ParseSemver(s string) error {
	if len(s) == 0 {
		return typex.Errorf(errors.Source, errors.InvalidArgument, "Semver string empty")
	}

	// Split into major.minor.(patch+pr+meta)
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return typex.Errorf(errors.Source, errors.InvalidArgument, "No Major.Minor.Patch elements found")
	}

	// Major
	if !containsOnly(parts[0], numbers) {
		return typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid character(s) found in major number %q", parts[0])
	}
	if hasLeadingZeroes(parts[0]) {
		return typex.Errorf(errors.Source, errors.InvalidArgument,
			"Major number must not contain leading zeroes %q", parts[0])
	}
	if _, err := strconv.ParseUint(parts[0], 10, 64); err != nil {
		return err
	}

	// Minor
	if !containsOnly(parts[1], numbers) {
		return typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid character(s) found in minor number %q", parts[1])
	}
	if hasLeadingZeroes(parts[1]) {
		return typex.Errorf(errors.Source, errors.InvalidArgument,
			"Minor number must not contain leading zeroes %q", parts[1])
	}
	if _, err := strconv.ParseUint(parts[1], 10, 64); err != nil {
		return err
	}

	patchStr := parts[2]
	if !containsOnly(patchStr, numbers) {
		return typex.Errorf(errors.Source, errors.InvalidArgument,
			"Invalid character(s) found in patch number %q", patchStr)
	}
	if hasLeadingZeroes(patchStr) {
		return typex.Errorf(errors.Source, errors.InvalidArgument,
			"Patch number must not contain leading zeroes %q", patchStr)
	}

	if _, err := strconv.ParseUint(patchStr, 10, 64); err != nil {
		return err
	}

	return nil
}

func containsOnly(s string, set string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(set, r)
	}) == -1
}

func hasLeadingZeroes(s string) bool {
	return len(s) > 1 && s[0] == '0'
}
