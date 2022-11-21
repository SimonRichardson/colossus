package mongo

import (
	"fmt"

	pool "github.com/SimonRichardson/colossus/fusion/mongo"
	"github.com/SimonRichardson/colossus/fusion/tests/stubs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func QueryFindOne(fn func(query bson.M) func(interface{}) error) pool.SessionCreator {
	return func(info *mgo.DialInfo) (pool.Session, error) {
		actions := stubs.NewActions()
		actions.
			On("DB", func(args ...interface{}) (interface{}, error) {
				return NewDatabase(actions), nil
			}).
			On("C", func(args ...interface{}) (interface{}, error) {
				return NewCollection(actions), nil
			}).
			On("Find", func(findArgs ...interface{}) (interface{}, error) {
				actions := stubs.NewActions()
				actions.
					On("One", func(oneArgs ...interface{}) (interface{}, error) {
						if m, ok := findArgs[0].(bson.M); ok {
							err := fn(m)(oneArgs[0])
							return nil, err
						}
						return nil, fmt.Errorf("Unexpected arguments")
					}).
					On("Sort", func(args ...interface{}) (interface{}, error) {
						return NewQuery(actions), nil
					}).
					On("Limit", func(args ...interface{}) (interface{}, error) {
						return NewQuery(actions), nil
					}).
					On("Skip", func(args ...interface{}) (interface{}, error) {
						return NewQuery(actions), nil
					}).
					On("Count", func(args ...interface{}) (interface{}, error) {
						return 0, nil
					})
				return NewQuery(actions), nil
			})

		return New(actions), nil
	}
}

func QueryFindAll(fn func(query bson.M) func(interface{}) error) pool.SessionCreator {
	return func(info *mgo.DialInfo) (pool.Session, error) {
		actions := stubs.NewActions()
		actions.
			On("DB", func(args ...interface{}) (interface{}, error) {
				return NewDatabase(actions), nil
			}).
			On("C", func(args ...interface{}) (interface{}, error) {
				return NewCollection(actions), nil
			}).
			On("Find", func(findArgs ...interface{}) (interface{}, error) {
				actions := stubs.NewActions()
				actions.
					On("All", func(allArgs ...interface{}) (interface{}, error) {
						if m, ok := findArgs[0].(bson.M); ok {
							err := fn(m)(allArgs[0])
							return nil, err
						}
						return nil, fmt.Errorf("Unexpected arguments")
					}).
					On("Sort", func(args ...interface{}) (interface{}, error) {
						return NewQuery(actions), nil
					}).
					On("Limit", func(args ...interface{}) (interface{}, error) {
						return NewQuery(actions), nil
					}).
					On("Skip", func(args ...interface{}) (interface{}, error) {
						return NewQuery(actions), nil
					})
				return NewQuery(actions), nil
			})

		return New(actions), nil
	}
}

func QueryFindAllAndCount(
	fn1 func(query bson.M) func(interface{}) error,
	fn2 func(query bson.M) func() (interface{}, error),
) pool.SessionCreator {
	return func(info *mgo.DialInfo) (pool.Session, error) {
		actions := stubs.NewActions()
		actions.
			On("DB", func(args ...interface{}) (interface{}, error) {
				return NewDatabase(actions), nil
			}).
			On("C", func(args ...interface{}) (interface{}, error) {
				return NewCollection(actions), nil
			}).
			On("Find", func(findArgs ...interface{}) (interface{}, error) {
				actions := stubs.NewActions()
				actions.
					On("All", func(allArgs ...interface{}) (interface{}, error) {
						if m, ok := findArgs[0].(bson.M); ok {
							err := fn1(m)(allArgs[0])
							return nil, err
						}
						return nil, fmt.Errorf("Unexpected arguments")
					}).
					On("Count", func(...interface{}) (interface{}, error) {
						if m, ok := findArgs[0].(bson.M); ok {
							return fn2(m)()
						}
						return nil, fmt.Errorf("Unexpected arguments")
					}).
					On("Sort", func(args ...interface{}) (interface{}, error) {
						return NewQuery(actions), nil
					}).
					On("Limit", func(args ...interface{}) (interface{}, error) {
						return NewQuery(actions), nil
					}).
					On("Skip", func(args ...interface{}) (interface{}, error) {
						return NewQuery(actions), nil
					})
				return NewQuery(actions), nil
			})

		return New(actions), nil
	}
}

func QueryPipeOne(fn func([]bson.D) func(interface{}) error) pool.SessionCreator {
	return func(info *mgo.DialInfo) (pool.Session, error) {
		actions := stubs.NewActions()
		actions.
			On("DB", func(args ...interface{}) (interface{}, error) {
				return NewDatabase(actions), nil
			}).
			On("C", func(args ...interface{}) (interface{}, error) {
				return NewCollection(actions), nil
			}).
			On("Pipe", func(findArgs ...interface{}) (interface{}, error) {
				actions := stubs.NewActions()
				actions.On("One", func(oneArgs ...interface{}) (interface{}, error) {
					if m, ok := findArgs[0].([]bson.D); ok {
						err := fn(m)(oneArgs[0])
						return nil, err
					}
					return nil, fmt.Errorf("Unexpected arguments")
				})
				return NewPipe(actions), nil
			})
		return New(actions), nil
	}
}

func QueryFindOneAndPipeOne(
	fn1 func(query bson.M) func(interface{}) error,
	fn2 func([]bson.D) func(interface{}) error,
) pool.SessionCreator {
	return func(info *mgo.DialInfo) (pool.Session, error) {
		actions := stubs.NewActions()
		actions.
			On("DB", func(args ...interface{}) (interface{}, error) {
				return NewDatabase(actions), nil
			}).
			On("C", func(args ...interface{}) (interface{}, error) {
				return NewCollection(actions), nil
			}).
			On("Find", func(findArgs ...interface{}) (interface{}, error) {
				actions := stubs.NewActions()
				actions.On("One", func(oneArgs ...interface{}) (interface{}, error) {
					if m, ok := findArgs[0].(bson.M); ok {
						err := fn1(m)(oneArgs[0])
						return nil, err
					}
					return nil, fmt.Errorf("Unexpected arguments")
				})
				return NewQuery(actions), nil
			}).
			On("Pipe", func(findArgs ...interface{}) (interface{}, error) {
				actions := stubs.NewActions()
				actions.On("One", func(oneArgs ...interface{}) (interface{}, error) {
					if m, ok := findArgs[0].([]bson.D); ok {
						err := fn2(m)(oneArgs[0])
						return nil, err
					}
					return nil, fmt.Errorf("Unexpected arguments")
				})
				return NewPipe(actions), nil
			})
		return New(actions), nil
	}
}

func QueryFindOneAndPipeOneAndCount(
	fn1 func(query bson.M) func(interface{}) error,
	fn2 func([]bson.D) func(interface{}) error,
	fn3 func(query bson.M) func() (interface{}, error),
) pool.SessionCreator {
	return func(info *mgo.DialInfo) (pool.Session, error) {
		actions := stubs.NewActions()
		actions.
			On("DB", func(args ...interface{}) (interface{}, error) {
				return NewDatabase(actions), nil
			}).
			On("C", func(args ...interface{}) (interface{}, error) {
				return NewCollection(actions), nil
			}).
			On("Find", func(findArgs ...interface{}) (interface{}, error) {
				actions := stubs.NewActions()
				actions.
					On("One", func(oneArgs ...interface{}) (interface{}, error) {
						if m, ok := findArgs[0].(bson.M); ok {
							err := fn1(m)(oneArgs[0])
							return nil, err
						}
						return nil, fmt.Errorf("Unexpected arguments")
					}).
					On("Count", func(...interface{}) (interface{}, error) {
						if m, ok := findArgs[0].(bson.M); ok {
							return fn3(m)()
						}
						return nil, fmt.Errorf("Unexpected arguments")
					})
				return NewQuery(actions), nil
			}).
			On("Pipe", func(findArgs ...interface{}) (interface{}, error) {
				actions := stubs.NewActions()
				actions.On("One", func(oneArgs ...interface{}) (interface{}, error) {
					if m, ok := findArgs[0].([]bson.D); ok {
						err := fn2(m)(oneArgs[0])
						return nil, err
					}
					return nil, fmt.Errorf("Unexpected arguments")
				})
				return NewPipe(actions), nil
			})
		return New(actions), nil
	}
}
