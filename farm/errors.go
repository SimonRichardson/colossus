package farm

import "fmt"

// PartialError defines an error case where expected things where meant to
// happen, but actually something different occured causing a partial update.
type PartialError interface {
	Expected() int
	Actual() []int
	Error() string
}

type PartialErrorType string

const (
	Store       PartialErrorType = "store"
	Counter     PartialErrorType = "counter"
	Persistence PartialErrorType = "persistence"
)

func NewPartialError(farm PartialErrorType, expected int, actual []int) PartialError {
	return partialError{farm, "Partial Error", expected, actual}
}

type partialError struct {
	farm     PartialErrorType
	reason   string
	expected int
	actual   []int
}

func (p partialError) Error() string {
	return fmt.Sprintf("%s from %s farm (Expected: %d, Actual: %v)",
		p.reason, p.farm, p.expected, p.actual,
	)
}

// Expected returns what partial values should have been seen
func (p partialError) Expected() int { return p.expected }

// Actual returns the values that where returned
func (p partialError) Actual() []int { return p.actual }

func PartialRepairError(err error, fn func()) error {
	if _, ok := err.(PartialError); ok {
		// Possibly dispatch a instrumentation call here!
		go func() { fn() }()
	}
	return err
}
