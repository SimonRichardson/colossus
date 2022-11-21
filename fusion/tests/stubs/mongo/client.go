package mongo

import (
	"gopkg.in/mgo.v2"

	pool "github.com/SimonRichardson/colossus/fusion/mongo"
	"github.com/SimonRichardson/colossus/fusion/tests/stubs"
	"github.com/SimonRichardson/colossus/typex"
)

type Session struct {
	*stubs.Actions
}

func New(actions *stubs.Actions) *Session {
	return &Session{actions}
}

func (s *Session) Ping() error        { return nil }
func (s *Session) Copy() pool.Session { return s }
func (s *Session) Close()             {}

func (s *Session) DB(name string) pool.Database {
	d, err := s.Run("DB", func(action stubs.Action) (interface{}, error) {
		return action.Run(name)
	})
	if err != nil {
		typex.Fatal(err)
	}
	return d.(pool.Database)
}

type Database struct {
	*stubs.Actions
}

func NewDatabase(actions *stubs.Actions) *Database {
	return &Database{actions}
}

func (d *Database) C(name string) pool.Collection {
	c, err := d.Run("C", func(action stubs.Action) (interface{}, error) {
		return action.Run(name)
	})
	if err != nil {
		typex.Fatal(err)
	}
	return c.(pool.Collection)
}

type Collection struct {
	*stubs.Actions
}

func NewCollection(actions *stubs.Actions) *Collection {
	return &Collection{actions}
}

func (c *Collection) Bulk() pool.Bulk                                            { return nil }
func (c *Collection) Insert(...interface{}) error                                { return nil }
func (c *Collection) UpsertId(interface{}, interface{}) (*mgo.ChangeInfo, error) { return nil, nil }
func (c *Collection) UpdateId(interface{}, interface{}) error                    { return nil }
func (c *Collection) Update(interface{}, interface{}) error                      { return nil }
func (c *Collection) Remove(interface{}) error                                   { return nil }
func (c *Collection) RemoveAll(interface{}) (*mgo.ChangeInfo, error)             { return nil, nil }

func (c *Collection) Pipe(query interface{}) pool.Pipe {
	p, err := c.Run("Pipe", func(action stubs.Action) (interface{}, error) {
		return action.Run(query)
	})
	if err != nil {
		typex.Fatal(err)
	}
	return p.(pool.Pipe)
}

func (c *Collection) Find(query interface{}) pool.Query {
	q, err := c.Run("Find", func(action stubs.Action) (interface{}, error) {
		return action.Run(query)
	})
	if err != nil {
		typex.Fatal(err)
	}
	return q.(pool.Query)
}

type Query struct {
	*stubs.Actions
}

func NewQuery(actions *stubs.Actions) *Query {
	return &Query{actions}
}

func (q *Query) One(result interface{}) error {
	res, err := q.Run("One", func(action stubs.Action) (interface{}, error) {
		return action.Run(result)
	})
	if err != nil {
		typex.Fatal(err)
	}
	switch t := res.(type) {
	case error:
		return t
	}
	return nil
}

func (q *Query) All(result interface{}) error {
	res, err := q.Run("All", func(action stubs.Action) (interface{}, error) {
		return action.Run(result)
	})
	if err != nil {
		typex.Fatal(err)
	}
	switch t := res.(type) {
	case error:
		return t
	}
	return nil
}

func (q *Query) Count() (int, error) {
	a, err := q.Run("Count", func(action stubs.Action) (interface{}, error) {
		return action.Run()
	})
	return a.(int), err
}

func (q *Query) Limit(limit int) pool.Query {
	a, err := q.Run("Limit", func(action stubs.Action) (interface{}, error) {
		return action.Run(limit)
	})
	if err != nil {
		typex.Fatal(err)
	}
	return a.(pool.Query)
}

func (q *Query) Skip(offset int) pool.Query {
	a, err := q.Run("Skip", func(action stubs.Action) (interface{}, error) {
		return action.Run(offset)
	})
	if err != nil {
		typex.Fatal(err)
	}
	return a.(pool.Query)
}

func (q *Query) Sort(values ...string) pool.Query {
	a, err := q.Run("Sort", func(action stubs.Action) (interface{}, error) {
		x := make([]interface{}, 0, len(values))
		for _, v := range values {
			x = append(x, v)
		}
		return action.Run(x...)
	})
	if err != nil {
		typex.Fatal(err)
	}
	return a.(pool.Query)
}

type Pipe struct {
	*stubs.Actions
}

func NewPipe(actions *stubs.Actions) *Pipe {
	return &Pipe{actions}
}

func (q *Pipe) All(result interface{}) error {
	res, err := q.Run("All", func(action stubs.Action) (interface{}, error) {
		return action.Run(result)
	})
	if err != nil {
		typex.Fatal(err)
	}
	switch t := res.(type) {
	case error:
		return t
	}
	return nil
}

func (q *Pipe) Explain(interface{}) error { return nil }

func (q *Pipe) One(result interface{}) error {
	res, err := q.Run("One", func(action stubs.Action) (interface{}, error) {
		return action.Run(result)
	})
	if err != nil {
		typex.Fatal(err)
	}
	switch t := res.(type) {
	case error:
		return t
	}
	return nil
}
