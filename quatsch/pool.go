package quatsch

type Pool interface {
	Sizable
	Get() (interface{}, error)
}

type Variance interface {
	Sizable
	With(func(interface{}) error) error
}

type Sizable interface {
	Len() int
}

type pool struct {
	variance Variance
}

func New(variance Variance) Pool {
	return &pool{variance}
}

func (p pool) Get() (res interface{}, err error) {
	err = p.variance.With(func(v interface{}) error {
		res = v
		return nil
	})
	return
}

func (p pool) Len() int {
	return p.variance.Len()
}
