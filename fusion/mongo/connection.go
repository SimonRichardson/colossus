package mongo

import (
	"strings"
	"sync"

	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/typex"

	"gopkg.in/mgo.v2"
)

type SessionCreator func(*mgo.DialInfo) (Session, error)

type ConnectionPool struct {
	mutex *sync.Mutex
	cond  *sync.Cond

	address string
	timeout *ConnectionTimeout
	creator SessionCreator

	primary     Session
	connections []Session
	loaned      int
	max         int
}

func NewConnectionPool(address string,
	timeout *ConnectionTimeout,
	maxConnections int,
	creator SessionCreator,
) *ConnectionPool {
	var (
		mutex        = &sync.Mutex{}
		primary, err = sessionCreator(creator)(&mgo.DialInfo{
			Addrs:   strings.Split(address, ","),
			Timeout: timeout.global,
		})
	)
	if err != nil {
		teleprinter.L.Error().Printf("Error connecting to Mongo Servers (addresses: %s) with error - %s\n",
			address, err.Error())
		typex.Fatal(err)
	}
	return &ConnectionPool{
		mutex: mutex,
		cond:  sync.NewCond(mutex),

		address: address,
		timeout: timeout,
		creator: sessionCreator(creator),

		primary:     primary,
		connections: []Session{},
		loaned:      0,
		max:         maxConnections,
	}
}

func (p *ConnectionPool) Get() (Session, error) {
	p.mutex.Lock()
	for {
		available := len(p.connections)
		switch {
		case available < 1 && p.loaned >= p.max:
			p.cond.Wait()

		case available < 1 && p.loaned < p.max:
			p.loaned++
			p.mutex.Unlock()

			return p.primary.Copy(), nil

		case available > 0:
			var session Session
			session, p.connections = p.connections[0], p.connections[1:]
			if p.loaned < p.max {
				p.loaned++
			}
			p.mutex.Unlock()
			return session.Copy(), nil
		}
	}
}

func (p *ConnectionPool) Put(session Session) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if session == nil || session.Ping() != nil {
		if p.loaned > 0 {
			p.loaned--
		}
		p.cond.Signal()
		return
	}

	p.connections = append(p.connections, session.Copy())
	go session.Close()

	if p.loaned > 0 {
		p.loaned--
	}
	p.cond.Signal()
}

func (p *ConnectionPool) CloseAll() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, conn := range p.connections {
		conn.Close()
	}
	p.connections = []Session{}
}

func sessionCreator(creator SessionCreator) (res SessionCreator) {
	if creator == nil {
		res = func(info *mgo.DialInfo) (Session, error) {
			conn, err := mgo.DialWithInfo(info)
			return &sess{conn}, err
		}
	} else {
		res = creator
	}
	return
}
