package redis

import (
	"fmt"

	"github.com/SimonRichardson/colossus/fusion"
	"github.com/SimonRichardson/colossus/fusion/errors"
	"github.com/SimonRichardson/colossus/typex"
	r "github.com/garyburd/redigo/redis"
)

// Pool is a collection of connections and a routing strategy.
type Pool struct {
	connections []*connectionPool
	routing     fusion.SelectionStrategy
}

// New creates a new Pool struct, which allows you to interact with the redis
// servers in away that is efficent as possible.
func New(addresses []string,
	routing fusion.SelectionStrategy,
	timeout *ConnectionTimeout,
	maxConnectionsPerInstance int,
	creator RedisCreator,
) *Pool {
	connections := make([]*connectionPool, len(addresses))
	for k, host := range addresses {
		uri, err := ParseRedisURL(host)
		if err != nil {
			// This is horrid, as can cause runtime errors.
			panic(fmt.Errorf("Invalid redis host %s : %v", host, err.Error()))
		}

		connections[k] = newConnectionPool(
			uri.String(),
			uri.Password(),
			timeout,
			maxConnectionsPerInstance,
			creator,
		)
	}
	return &Pool{connections, routing}
}

// Size returns the number of connections that the pool is holding on to.
func (p *Pool) Size() int {
	return len(p.connections)
}

// Index returns the pool connection index depending on the key and the number
// of connections the pool has. Note: that it can ignore the key depending on
// the routing strategy.
func (p *Pool) Index(key string) int {
	size := p.Size()
	if size <= 1 {
		return 0
	}
	return p.routing.Select(key, size)
}

// With defines a function that will execute the function argument when it
// locates a pool connection.
func (p *Pool) With(key string, do func(r.Conn) error) error {
	return p.WithIndex(p.Index(key), do)
}

// WithIndex defines a way to target a pool connection directly. If you already
// know the index this is the most efficent way to get a connection. There is
// no bounds checking, so runtime errors can occur!
func (p *Pool) WithIndex(index int, do func(r.Conn) error) error {
	if len(p.connections) < index {
		return typex.Errorf(errors.Source, errors.UnexpectedArgument, "Invalid index")
	}

	conn, err := p.connections[index].get()
	defer p.connections[index].put(conn)
	if err != nil {
		return err
	}

	return do(conn)
}

// Close closes all available connections with in the pool.
func (p *Pool) Close() {
	for _, conn := range p.connections {
		conn.closeAll()
	}
	p.connections = []*connectionPool{}
}
