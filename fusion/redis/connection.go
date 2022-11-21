package redis

import (
	"sync"

	"github.com/SimonRichardson/colossus/teleprinter"
	r "github.com/garyburd/redigo/redis"
)

type RedisCreator func(string, string, *ConnectionTimeout) (r.Conn, error)

type connectionPool struct {
	mutex *sync.Mutex
	cond  *sync.Cond

	address, password string

	timeout *ConnectionTimeout
	creator RedisCreator

	connections []r.Conn
	loaned      int
	max         int
}

func newConnectionPool(address, password string,
	timeout *ConnectionTimeout,
	maxConnections int,
	creator RedisCreator,
) *connectionPool {
	mutex := &sync.Mutex{}
	return &connectionPool{
		mutex: mutex,
		cond:  sync.NewCond(mutex),

		address:  address,
		password: password,

		timeout: timeout,
		creator: redisCreator(creator),

		connections: []r.Conn{},
		loaned:      0,
		max:         maxConnections,
	}
}

func (p *connectionPool) get() (r.Conn, error) {
	p.mutex.Lock()
start_loop:
	for {
		available := len(p.connections)
		switch {
		case available < 1 && p.loaned >= p.max:
			p.cond.Wait()

		case available < 1 && p.loaned < p.max:
			p.loaned++
			p.mutex.Unlock()
			return p.creator(p.address, p.password, p.timeout)

		case available > 0:
			var conn r.Conn
			conn, p.connections = p.connections[0], p.connections[1:]
			// Recusive request - note this will drain the pool if there are
			// errors on the connection
			if conn == nil {
				continue start_loop
			}
			if conn.Err() != nil {
				conn.Close()
				continue start_loop
			}

			if p.loaned < p.max {
				p.loaned++
			}

			p.mutex.Unlock()
			return conn, nil
		}
	}
}

func (p *connectionPool) put(conn r.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if conn == nil || conn.Err() != nil {
		if p.loaned > 0 {
			p.loaned--
		}
		// Make sure we close all connections.
		if conn != nil && conn.Err() != nil {
			conn.Close()
		}
		p.cond.Signal()
		return
	}

	if len(p.connections) >= p.max {
		conn.Close()
		return
	}

	p.connections = append(p.connections, conn)
	if p.loaned > 0 {
		p.loaned--
	}
	p.cond.Signal()
}

func (p *connectionPool) closeAll() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, conn := range p.connections {
		conn.Close()
	}
	p.connections = []r.Conn{}
}

func redisCreator(creator RedisCreator) (res RedisCreator) {
	if creator == nil {
		res = func(address, password string, timeout *ConnectionTimeout) (r.Conn, error) {
			conn, err := r.Dial("tcp", address,
				r.DialPassword(password),
				r.DialConnectTimeout(timeout.connect),
				r.DialReadTimeout(timeout.read),
				r.DialWriteTimeout(timeout.write),
			)
			if err != nil {
				teleprinter.L.Error().Printf("Error connecting to Redis Servers (addresses: %s) with error - %s\n",
					address, err.Error())
			}
			return conn, err
		}
	} else {
		res = creator
	}
	return
}
