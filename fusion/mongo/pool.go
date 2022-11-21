package mongo

import "github.com/SimonRichardson/colossus/fusion"

// Pool is a collection of connections and a routing strategy.
type Pool struct {
	connections []*ConnectionPool
	routing     fusion.SelectionStrategy
}

// New creates a new Pool struct, which allows you to interact with the redis
// servers in away that is efficent as possible.
func New(addresses []string,
	routing fusion.SelectionStrategy,
	timeout *ConnectionTimeout,
	maxConnectionsPerInstance int,
	creator SessionCreator,
) *Pool {
	connections := make([]*ConnectionPool, len(addresses))
	for k, address := range addresses {
		connections[k] = NewConnectionPool(address,
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
func (p *Pool) With(key string, do func(Session) error) error {
	return p.WithIndex(p.Index(key), do)
}

// WithIndex defines a way to target a pool connection directly. If you already
// know the index this is the most efficent way to get a connection. There is
// no bounds checking, so runtime errors can occur!
func (p *Pool) WithIndex(index int, do func(Session) error) error {
	session, err := p.connections[index].Get()
	defer p.connections[index].Put(session)
	if err != nil {
		return err
	}

	return do(session)
}

// Close closes all available connections with in the pool.
func (p *Pool) Close() {
	for _, pool := range p.connections {
		pool.CloseAll()
	}
}
