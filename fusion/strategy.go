package fusion

// SelectionStrategy defines an interface for selecting which pool to use
// depending on the key and how many connections there is.
type SelectionStrategy interface {
	Select(string, int) int
}
