package notifier

import (
	"sync"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	t "github.com/SimonRichardson/colossus/cluster"
	p "github.com/SimonRichardson/colossus/fusion/redis"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/garyburd/redigo/redis"
)

const (
	defaultSubscribeKey  string = "P5u0fIX6BJmM1OIOi7L8qnjDIZcOKwOtM"
	defaultVerifyResults bool   = true
)

// Cluster defines an interface for accessing redis storage. It defines various
// strategies for querying the redis database and the most efficient way
// possible.
type Cluster interface {
	t.Notifier
	t.Closer
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

func (c *cluster) Publish(channel s.Channel, members []s.KeyFieldScoreSizeExpiry) <-chan t.Element {
	keys, values := s.KeyFieldScoreSizeExpiries(members).Bucketize()
	return c.common(keys, func(conn redis.Conn, key bs.Key) error {
		return publish(conn, channel, values[key])
	})
}

func (c *cluster) Unpublish(channel s.Channel, members []s.KeyFieldScoreSizeExpiry) <-chan t.Element {
	keys, values := s.KeyFieldScoreSizeExpiries(members).Bucketize()
	return c.common(keys, func(conn redis.Conn, key bs.Key) error {
		return unpublish(conn, channel, values[key])
	})
}

func (c *cluster) Subscribe(channel s.Channel) <-chan t.Element {
	return c.subscribe(func(conn redis.Conn) <-chan t.Element {
		return subscribe(conn, channel)
	})
}

func (c *cluster) Close() error {
	c.pool.Close()
	return nil
}

func (c *cluster) common(keys []bs.Key, f func(redis.Conn, bs.Key) error) <-chan t.Element {
	out := make(chan t.Element)
	go func() {

		wg := sync.WaitGroup{}
		wg.Add(len(keys))

		for k, v := range keys {
			go func(index int, key bs.Key) {
				defer wg.Done()

				var elements []t.Element
				if err := c.pool.With(key.String(), func(conn redis.Conn) (err error) {
					err = f(conn, key)
					return
				}); err != nil {
					elements = errorElementsFromKeyCount(key, err)
				}

				for _, element := range elements {
					out <- element
				}
			}(k, v)
		}

		wg.Wait()
		close(out)
	}()
	return out
}

func errorElementsFromKeyCount(key bs.Key, err error) []t.Element {
	return []t.Element{t.NewErrorElement(key, err)}
}

func (c *cluster) subscribe(f func(redis.Conn) <-chan t.Element) <-chan t.Element {
	out := make(chan t.Element)
	go func() {
		if err := c.pool.With(defaultSubscribeKey, func(conn redis.Conn) error {
			go func() {
				for element := range f(conn) {
					out <- element
				}
			}()
			return nil
		}); err != nil {
			out <- t.NewErrorElement(bs.Key(defaultSubscribeKey), err)
		}
	}()
	return out
}
