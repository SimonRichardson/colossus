package counter

import (
	"sync"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	t "github.com/SimonRichardson/colossus/cluster"
	p "github.com/SimonRichardson/colossus/fusion/redis"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/garyburd/redigo/redis"
)

// DefaultKeysKey is a random key so that we can query the pool with a key and
// we want a constant one, as the query actually scans the whole thing, so no
// identifier is identifiable.
const (
	defaultKeysKey   string = "IiDIZq07L8nOKwP6BjfO5ucJmM1OIOXtM"
	defaultBatchSize int    = 1000

	defaultVerifyResults bool = true
)

// Cluster defines an interface for accessing redis storage. It defines various
// strategies for querying the redis database and the most efficient way
// possible.
type Cluster interface {
	t.Inserter
	t.Deleter
	t.Scanner
	t.Scorer
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

func (c *cluster) Insert(members []s.KeyFieldScoreTxnValue, sizeExpiry s.KeySizeExpiry) <-chan t.Element {
	keys, values := s.KeyFieldScoreTxnValues(members).KeysBucketize()
	return c.countCommon(keys, func(conn redis.Conn, key bs.Key) ([]s.KeyCount, error) {
		return insertion(conn, values[key], sizeExpiry[key])
	})
}

func (c *cluster) Delete(members []s.KeyFieldScoreTxnValue, sizeExpiry s.KeySizeExpiry) <-chan t.Element {
	keys, values := s.KeyFieldScoreTxnValues(members).KeysBucketize()
	return c.countCommon(keys, func(conn redis.Conn, key bs.Key) ([]s.KeyCount, error) {
		return deletion(conn, values[key], sizeExpiry[key])
	})
}

func (c *cluster) Rollback(members []s.KeyFieldScoreTxnValue, sizeExpiry s.KeySizeExpiry) <-chan t.Element {
	return c.Delete(members, sizeExpiry)
}

func (c *cluster) Size(key bs.Key) <-chan t.Element {
	return c.countCommon([]bs.Key{key}, func(conn redis.Conn, key bs.Key) ([]s.KeyCount, error) {
		return cardinality(conn, key)
	})
}

func (c *cluster) Keys() <-chan t.Element {
	return c.keyCommon(bs.Key(defaultKeysKey), func(conn redis.Conn) ([]bs.Key, error) {
		return keys(conn, defaultBatchSize)
	})
}

func (c *cluster) Members(key bs.Key) <-chan t.Element {
	return c.keyCommon(key, func(conn redis.Conn) ([]bs.Key, error) {
		return members(conn, key)
	})
}

func (c *cluster) Score(members []s.KeyFieldTxnValue) (map[s.KeyFieldTxnValue]s.Presence, error) {
	return c.scoreCommon(members, func(conn redis.Conn, members []s.KeyFieldTxnValue) (map[s.KeyFieldTxnValue]s.Presence, error) {
		return score(conn, members)
	})
}

func (c *cluster) Close() error {
	c.pool.Close()
	return nil
}

func (c *cluster) countCommon(keys []bs.Key, f func(redis.Conn, bs.Key) ([]s.KeyCount, error)) <-chan t.Element {
	out := make(chan t.Element)
	go func() {

		wg := sync.WaitGroup{}
		wg.Add(len(keys))

		for k, v := range keys {
			go func(index int, key bs.Key) {
				defer wg.Done()

				var (
					elements []t.Element
					result   = []s.KeyCount{s.KeyCount{Key: key, Count: 0}}
				)
				if err := c.pool.With(key.String(), func(conn redis.Conn) (err error) {
					result, err = f(conn, key)
					return
				}); err != nil {
					elements = errorElementsFromKeyCount(result, err)
				} else {
					elements = successElementsFromKeyCount(result)
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

func errorElementsFromKeyCount(keys []s.KeyCount, err error) []t.Element {
	elements := make([]t.Element, 0, len(keys))
	for _, k := range keys {
		elements = append(elements, t.NewErrorElement(k.Key, err))
	}
	return elements
}

func successElementsFromKeyCount(keys []s.KeyCount) []t.Element {
	buckets := map[bs.Key][]s.KeyCount{}
	for _, v := range keys {
		buckets[v.Key] = append(buckets[v.Key], v)
	}

	elements := make([]t.Element, 0, len(keys))
	for k, v := range buckets {
		count := 0
		for _, c := range v {
			count += c.Count
		}
		elements = append(elements, t.NewCountElement(k, count))
	}

	return elements
}

func (c *cluster) keyCommon(key bs.Key,
	f func(redis.Conn) ([]bs.Key, error),
) <-chan t.Element {
	out := make(chan t.Element)
	go func() {

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()

			var result []bs.Key
			if err := c.pool.With(key.String(), func(conn redis.Conn) (err error) {
				result, err = f(conn)
				return
			}); err != nil {
				out <- t.NewErrorElement(key, err)
			} else {
				out <- t.NewKeyElement(key, result)
			}
		}()

		wg.Wait()

		close(out)
	}()
	return out
}

func (c *cluster) scoreCommon(members []s.KeyFieldTxnValue,
	f func(redis.Conn, []s.KeyFieldTxnValue) (map[s.KeyFieldTxnValue]s.Presence, error),
) (map[s.KeyFieldTxnValue]s.Presence, error) {

	m := map[int][]s.KeyFieldTxnValue{}
	for _, keyField := range members {
		index := c.pool.Index(keyField.Key.String())
		m[index] = append(m[index], keyField)
	}

	type response struct {
		presenceMap map[s.KeyFieldTxnValue]s.Presence
		err         error
	}

	out := make(chan response, len(m))
	for index, keyFields := range m {
		go func(index int, keyFields []s.KeyFieldTxnValue) {
			var presenceMap map[s.KeyFieldTxnValue]s.Presence
			err := c.pool.WithIndex(index, func(conn redis.Conn) (err error) {
				presenceMap, err = f(conn, keyFields)
				return
			})
			out <- response{presenceMap, err}
		}(index, keyFields)
	}

	presenceMap := map[s.KeyFieldTxnValue]s.Presence{}
	for i := 0; i < cap(out); i++ {
		response := <-out
		if response.err != nil {
			continue
		}

		for keyField, presence := range response.presenceMap {
			presenceMap[keyField] = presence
		}
	}
	return presenceMap, nil
}
