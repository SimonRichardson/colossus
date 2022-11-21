package cache

import (
	"sync"
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/fusion/redis"
	r "github.com/garyburd/redigo/redis"
)

const (
	defaultKey = "uE7Lz8Swq0iDI51XtO6cfnOM1DIOMO8IKZPBjJmO"
)

// Cluster defines a interface for parallelised requesting.
type Cluster interface {
	sv.Encoder
}

type cluster struct {
	pool   *redis.Pool
	expiry time.Duration
}

func newCluster(pool *redis.Pool, expiry time.Duration) *cluster {
	return &cluster{pool, expiry}
}

func (c *cluster) GetBytes(key bs.Key) <-chan sv.Element {
	return c.common(func(conn r.Conn, dst chan sv.Element) {
		res, err := r.Bytes(conn.Do("GET", key.String()))
		if err != nil {
			dst <- sv.NewErrorElement(err)
			return
		}
		dst <- sv.NewBytesElement(res)
	})
}

func (c *cluster) SetBytes(key bs.Key, bytes []byte) <-chan sv.Element {
	return c.common(func(conn r.Conn, dst chan sv.Element) {
		_, err := conn.Do("SETEX", key.String(), c.expiry.Seconds(), bytes)
		dst <- sv.NewErrorElement(err)
	})
}

func (c *cluster) DelBytes(key bs.Key) <-chan sv.Element {
	return c.common(func(conn r.Conn, dst chan sv.Element) {
		_, err := conn.Do("DEL", key.String())
		dst <- sv.NewErrorElement(err)
	})
}

func (c *cluster) common(fn func(r.Conn, chan sv.Element)) <-chan sv.Element {
	out := make(chan sv.Element)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := c.pool.With(defaultKey, func(conn r.Conn) error {
				fn(conn, out)
				return nil
			}); err != nil {
				out <- sv.NewErrorElement(err)
			}
		}()

		wg.Wait()
		close(out)
	}()
	return out
}
