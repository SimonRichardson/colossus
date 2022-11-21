package score

import (
	"sync"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	t "github.com/SimonRichardson/colossus/cluster"
	u "github.com/SimonRichardson/colossus/colossus-shim/cluster"
	p "github.com/SimonRichardson/colossus/fusion/redis"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/garyburd/redigo/redis"
)

// Cluster defines an interface for accessing redis storage. It defines various
// strategies for querying the redis database and the most efficient way
// possible.
type Cluster interface {
	u.Incrementer
}

type cluster struct {
	pool *p.Pool
}

// New creates a cluster using a pool and ops to query the redis storage.
func New(pool *p.Pool) Cluster {
	return &cluster{
		pool: pool,
	}
}

func (c *cluster) Increment(key bs.Key) <-chan t.Element {
	return c.countCommon(key, func(conn redis.Conn, key bs.Key) (s.KeyCount, error) {
		return increment(conn, key)
	})
}

func (c *cluster) Close() error {
	c.pool.Close()
	return nil
}

func (c *cluster) countCommon(key bs.Key, f func(redis.Conn, bs.Key) (s.KeyCount, error)) <-chan t.Element {
	out := make(chan t.Element)
	go func() {

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func(key bs.Key) {
			defer wg.Done()

			result := s.KeyCount{Key: key, Count: 0}
			if err := c.pool.With(key.String(), func(conn redis.Conn) (err error) {
				result, err = f(conn, key)
				return
			}); err != nil {
				out <- t.NewErrorElement(key, err)
			} else {
				out <- t.NewCountElement(key, result.Count)
			}
		}(key)

		wg.Wait()
		close(out)
	}()
	return out
}
