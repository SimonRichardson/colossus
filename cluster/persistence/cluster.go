package persistence

import (
	"fmt"
	"sync"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	t "github.com/SimonRichardson/colossus/cluster"
	p "github.com/SimonRichardson/colossus/fusion/mongo"
	s "github.com/SimonRichardson/colossus/selectors"
)

// Cluster defines an interface for accessing redis storage. It defines various
// strategies for querying the redis database and the most efficient way
// possible.
type Cluster interface {
	t.Inserter
	t.Deleter
	t.Repairer
	t.Closer
}

type cluster struct {
	pool        *p.Pool
	dbName      string
	keyPrefix   string
	transformer s.Transformer
}

// New creates a cluster using a pool and ops to query the redis storage.
func New(pool *p.Pool,
	dbName, keyPrefix string,
	transformer s.Transformer,
) Cluster {
	return &cluster{
		pool:        pool,
		dbName:      dbName,
		keyPrefix:   keyPrefix,
		transformer: transformer,
	}
}

func (c *cluster) Insert(members []s.KeyFieldScoreTxnValue, sizeExpiry s.KeySizeExpiry) <-chan t.Element {
	keys, values := s.KeyFieldScoreTxnValues(members).KeysBucketize()
	return c.countCommon(keys, func(db p.Database, key bs.Key) ([]s.KeyCount, error) {
		return insertion(db, c.transformer, values[key])
	})
}

func (c *cluster) Delete(members []s.KeyFieldScoreTxnValue, sizeExpiry s.KeySizeExpiry) <-chan t.Element {
	keys, values := s.KeyFieldScoreTxnValues(members).KeysBucketize()
	return c.countCommon(keys, func(db p.Database, key bs.Key) ([]s.KeyCount, error) {
		return deletion(db, c.transformer, values[key])
	})
}

func (c *cluster) Rollback(members []s.KeyFieldScoreTxnValue, sizeExpiry s.KeySizeExpiry) <-chan t.Element {
	return c.Delete(members, sizeExpiry)
}

func (c *cluster) Repair(members []s.KeyFieldScoreTxnValue, sizeExpiry s.KeySizeExpiry) <-chan t.Element {
	keys, values := s.KeyFieldScoreTxnValues(members).KeysBucketize()
	return c.countCommon(keys, func(db p.Database, key bs.Key) ([]s.KeyCount, error) {
		return repair(db, c.transformer, values[key])
	})
}

func (c *cluster) Close() error {
	c.pool.Close()
	return nil
}

func (c *cluster) countCommon(keys []bs.Key, f func(p.Database, bs.Key) ([]s.KeyCount, error)) <-chan t.Element {
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
				if err := c.pool.With(key.String(), func(sess p.Session) (err error) {
					result, err = f(c.database(sess), key)
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

func (c *cluster) database(sess p.Session) p.Database {
	return sess.DB(c.dbName)
}

func (c *cluster) namespace(key bs.Key) bs.Key {
	return bs.Key(fmt.Sprintf("%s%s", c.keyPrefix, key.String()))
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
