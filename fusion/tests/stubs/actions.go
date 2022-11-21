package stubs

import "fmt"

type Action interface {
	Name() string
	Run(...interface{}) (interface{}, error)
	Once() bool
}

func MakeAction(name string, fn func(...interface{}) (interface{}, error)) Action {
	return action{name, fn, false}
}

type action struct {
	name string
	fn   func(...interface{}) (interface{}, error)
	once bool
}

func (s action) Name() string                                 { return s.name }
func (s action) Run(args ...interface{}) (interface{}, error) { return s.fn(args...) }
func (s action) Once() bool                                   { return s.once }

type Actions struct {
	Actions []Action
}

func NewActions() *Actions {
	return &Actions{[]Action{}}
}

func (t *Actions) Concat(b *Actions) *Actions {
	for _, v := range b.Actions {
		t.Actions = append(t.Actions, v)
	}
	return t
}

func (t *Actions) On(name string, fn func(...interface{}) (interface{}, error)) *Actions {
	t.Actions = append(t.Actions, MakeAction(name, fn))
	return t
}

func (t *Actions) Run(name string, fn func(Action) (interface{}, error)) (interface{}, error) {
	for k, v := range t.Actions {
		if v.Name() == name {
			res, err := fn(v)

			if v.Once() {
				t.Actions = append(t.Actions[:k], t.Actions[k+1:]...)
			}
			return res, err
		}
	}
	panic(fmt.Errorf("No Action found for name: %q.", name))
}
