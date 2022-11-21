package mongo

import (
	"github.com/SimonRichardson/colossus/fusion/errors"
	"github.com/SimonRichardson/colossus/typex"
	"gopkg.in/mgo.v2"
)

type Session interface {
	Ping() error
	Copy() Session
	Close()
	DB(string) Database
}

type Database interface {
	C(string) Collection
}

type Collection interface {
	Bulk() Bulk
	Pipe(interface{}) Pipe
	Find(interface{}) Query
	Insert(...interface{}) error
	UpsertId(interface{}, interface{}) (*mgo.ChangeInfo, error)
	UpdateId(interface{}, interface{}) error
	Update(interface{}, interface{}) error
	Remove(interface{}) error
	RemoveAll(interface{}) (*mgo.ChangeInfo, error)
}

type Bulk interface {
	Insert(...interface{})
	Unordered()
	Run() (*mgo.BulkResult, error)
}

type Pipe interface {
	One(interface{}) error
	All(interface{}) error
	Explain(interface{}) error
}

type Query interface {
	One(interface{}) error
	All(interface{}) error
	Count() (int, error)
	Limit(int) Query
	Skip(int) Query
	Sort(...string) Query
}

type sess struct {
	session *mgo.Session
}

func (s *sess) Ping() error {
	return s.with(func(s *mgo.Session) error {
		return s.Ping()
	})
}

func (s *sess) Copy() Session {
	return &sess{s.session.Copy()}
}

func (s *sess) Close() {
	s.with(func(s *mgo.Session) error {
		s.Close()
		return nil
	})
}

func (s *sess) DB(name string) Database {
	return &db{s.session.DB(name)}
}

func (s *sess) with(fn func(*mgo.Session) error) (err error) {
	if s == nil || s.session == nil {
		return typex.Errorf(errors.Source, errors.UnexpectedArgument, "Session was nil")
	}

	// This looks confusing, but is to prevent mgo panic'ing because it can't
	// deal with errors correctly!
	defer func() {
		switch e := recover().(type) {
		case nil:
			return
		case error:
			err = e
		default:
			typex.PrintStack(false)
			panic(e)
		}
	}()

	return fn(s.session)
}

type db struct {
	database *mgo.Database
}

func (d *db) C(name string) Collection {
	return &col{d.database.C(name)}
}

type col struct {
	collection *mgo.Collection
}

func (c *col) Bulk() Bulk {
	return &bulk{c.collection.Bulk()}
}

func (c *col) Pipe(pipeline interface{}) Pipe {
	return &pipe{c.collection.Pipe(pipeline)}
}

func (c *col) Find(selector interface{}) Query {
	return &query{c.collection.Find(selector)}
}

func (c *col) Insert(docs ...interface{}) error {
	return c.collection.Insert(docs...)
}

func (c *col) UpsertId(id interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	return c.collection.UpsertId(id, update)
}

func (c *col) UpdateId(id interface{}, update interface{}) error {
	return c.collection.UpdateId(id, update)
}

func (c *col) Update(selector interface{}, update interface{}) error {
	return c.collection.Update(selector, update)
}

func (c *col) Remove(selector interface{}) error {
	return c.collection.Remove(selector)
}

func (c *col) RemoveAll(selector interface{}) (*mgo.ChangeInfo, error) {
	return c.collection.RemoveAll(selector)
}

type bulk struct {
	bulk *mgo.Bulk
}

func (b *bulk) Insert(docs ...interface{}) {
	b.bulk.Insert(docs...)
}

func (b *bulk) Unordered() {
	b.bulk.Unordered()
}

func (b *bulk) Run() (*mgo.BulkResult, error) {
	return b.bulk.Run()
}

type pipe struct {
	pipe *mgo.Pipe
}

func (p *pipe) One(result interface{}) error {
	return p.pipe.One(result)
}

func (p *pipe) All(result interface{}) error {
	return p.pipe.All(result)
}

func (p *pipe) Explain(result interface{}) error {
	return p.pipe.Explain(result)
}

type query struct {
	query *mgo.Query
}

func (q *query) One(result interface{}) error {
	return q.query.One(result)
}

func (q *query) All(result interface{}) error {
	return q.query.All(result)
}

func (q *query) Count() (int, error) {
	return q.query.Count()
}

func (q *query) Limit(n int) Query {
	return &query{q.query.Limit(n)}
}

func (q *query) Skip(n int) Query {
	return &query{q.query.Skip(n)}
}

func (q *query) Sort(fields ...string) Query {
	return &query{q.query.Sort(fields...)}
}
