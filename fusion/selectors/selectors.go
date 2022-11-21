package selectors

type KeyStore interface {
	List(Prefix) (map[string]int, error)
}
