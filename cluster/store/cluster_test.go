package store

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"testing/quick"
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	c "github.com/SimonRichardson/colossus/cluster"
	"github.com/SimonRichardson/colossus/env"
	p "github.com/SimonRichardson/colossus/fusion/redis"
	"github.com/SimonRichardson/colossus/fusion/strategies"
	stubs "github.com/SimonRichardson/colossus/fusion/tests/stubs/redis"
	"github.com/SimonRichardson/colossus/quatsch"
	b "github.com/SimonRichardson/colossus/quatsch/pool/bson"
	"github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/tests"
	"github.com/SimonRichardson/colossus/typex"
)

var (
	defaultUseStubs = false
)

func TestMain(t *testing.M) {
	var flagStubs bool
	flag.BoolVar(&flagStubs, "stubs", false, "enable stubs testing")
	flag.Parse()

	defaultUseStubs = flagStubs

	os.Exit(t.Run())
}

func newCluster(creator p.RedisCreator) Cluster {
	var (
		e         = env.New(nil)
		clusters  = strings.Split(e.StoreInstances, ";")
		instances = strings.Split(clusters[0], ",")[0:1]

		pool = p.New(instances, strategies.NewHash(), &p.ConnectionTimeout{}, 100, creator)
	)

	return New(pool)
}

type fnAlias func(string, string, string, string, time.Duration) <-chan c.Element

func execute(fn func([]selectors.KeyFieldScoreTxnValue, selectors.KeySizeExpiry) <-chan c.Element, amount int) fnAlias {
	return func(key, field, txn, value string, duration time.Duration) <-chan c.Element {

		members := []selectors.KeyFieldScoreTxnValue{}
		for i := 0; i < amount; i++ {
			members = append(members, selectors.KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(fmt.Sprintf("%s_%d", field, i)),
				Score: 1,
				Txn:   bs.Key(txn),
				Value: value,
			})
		}
		return fn(members, selectors.KeySizeExpiry{
			bs.Key(key): selectors.SizeExpiry{
				Size:   int64(amount) + 1,
				Expiry: duration,
			},
		})
	}
}

func insert(cluster Cluster, amount int) fnAlias {
	return execute(cluster.Insert, amount)
}

func delete(cluster Cluster, amount int) fnAlias {
	return execute(cluster.Delete, amount)
}

func checkErrors(e <-chan c.Element) {
	for v := range e {
		if err := c.ErrorFromElement(v); err != nil {
			typex.Fatal(err)
		}
	}
}

func getIdentPool() quatsch.Pool {
	var (
		maxBuffer               = 99999
		maxInsertionPerDuration = int64(1000000)
	)

	return quatsch.New(b.New(maxBuffer, time.Second, maxInsertionPerDuration))
}

func TestInsert(t *testing.T) {
	var creator p.RedisCreator
	if defaultUseStubs {
		creator = stubs.SendAndReceive(
			func(name string, args ...interface{}) error {
				return nil
			},
			func() (interface{}, error) {
				return int64(1), nil
			},
		)
	}

	var (
		amount  = rand.Intn(5) + 1
		cluster = newCluster(creator)
		in      = insert(cluster, amount)
		pool    = getIdentPool()

		f = func(field, txn, value string, duration time.Duration) bool {
			key, err := b.Bson(pool.Get())
			if err != nil {
				typex.Fatal(err)
			}

			dst := in(key.Hex(), field, txn, value, duration)

			result := 0
			for e := range dst {
				if err := c.ErrorFromElement(e); err != nil {
					typex.Fatal(err)
				}
				result += c.AmountFromElement(e)
			}
			return result == amount
		}
	)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	var creator p.RedisCreator
	if defaultUseStubs {
		creator = stubs.SendAndReceive(
			func(name string, args ...interface{}) error {
				return nil
			},
			func() (interface{}, error) {
				return int64(1), nil
			},
		)
	}

	var (
		amount  = rand.Intn(5) + 1
		cluster = newCluster(creator)
		in      = insert(cluster, amount)
		del     = delete(cluster, amount)
		pool    = getIdentPool()

		f = func(field, txn, value string, duration time.Duration) bool {
			key, err := b.Bson(pool.Get())
			if err != nil {
				typex.Fatal(err)
			}

			in(key.Hex(), field, txn, value, duration)
			dst := del(key.Hex(), field, txn, value, duration)

			result := 0
			for e := range dst {
				if err := c.ErrorFromElement(e); err != nil {
					typex.Fatal(err)
				}
				result += c.AmountFromElement(e)
			}
			return result == amount
		}
	)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

func TestSize(t *testing.T) {
	var (
		amount  = rand.Intn(5) + 1
		creator p.RedisCreator
	)
	if defaultUseStubs {
		creator = stubs.DoSendAndReceive(
			func(name string, args ...interface{}) (interface{}, error) {
				return int64(amount), nil
			},
			func(name string, args ...interface{}) error {
				return nil
			},
			func() (interface{}, error) {
				return int64(1), nil
			},
		)
	}

	var (
		cluster = newCluster(creator)
		in      = insert(cluster, amount)
		pool    = getIdentPool()

		f = func(field, txn, value string, duration time.Duration) bool {
			key, err := b.Bson(pool.Get())
			if err != nil {
				typex.Fatal(err)
			}

			checkErrors(in(key.Hex(), field, txn, value, duration))
			dst := cluster.Size(bs.Key(key.Hex()))

			result := 0
			for e := range dst {
				if err := c.ErrorFromElement(e); err != nil {
					typex.Fatal(err)
				}
				result += c.AmountFromElement(e)
			}
			return result == amount
		}
	)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

func TestKeys(t *testing.T) {
	var (
		amount = rand.Intn(5) + 1
		keys   = []interface{}{}

		creator p.RedisCreator
	)
	if defaultUseStubs {
		creator = stubs.DoSendAndReceive(
			func(name string, args ...interface{}) (interface{}, error) {
				if name == "SCAN" {
					return []interface{}{int64(0), keys}, nil
				}
				return nil, nil
			},
			func(name string, args ...interface{}) error {
				return nil
			},
			func() (interface{}, error) {
				return int64(1), nil
			},
		)
	}

	var (
		cluster = newCluster(creator)
		in      = insert(cluster, amount)
		pool    = getIdentPool()

		f = func(field, txn, value string, duration time.Duration) bool {
			key, err := b.Bson(pool.Get())
			if err != nil {
				typex.Fatal(err)
			}
			val := key.Hex()
			keys = append(keys, []byte(fmt.Sprintf("%s%s%s", prefix, val, insertSuffix)))

			checkErrors(in(val, field, txn, value, duration))
			dst := cluster.Keys()

			for e := range dst {
				if err := c.ErrorFromElement(e); err != nil {
					typex.Fatal(err)
				}
				for _, v := range c.KeysFromElement(e) {
					if v.String() == val {
						return true
					}
				}
			}
			return false
		}
	)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

func TestMembers(t *testing.T) {
	var (
		amount = rand.Intn(5) + 1
		fields = []interface{}{}

		creator p.RedisCreator
	)
	if defaultUseStubs {
		creator = stubs.DoSendAndReceive(
			func(name string, args ...interface{}) (interface{}, error) {
				if name == "HKEYS" {
					return fields, nil
				}
				return nil, nil
			},
			func(name string, args ...interface{}) error {
				return nil
			},
			func() (interface{}, error) {
				return int64(1), nil
			},
		)
	}

	var (
		cluster = newCluster(creator)
		in      = insert(cluster, amount)
		pool    = getIdentPool()

		f = func(field, txn, value string, duration time.Duration) bool {
			key, err := b.Bson(pool.Get())
			if err != nil {
				typex.Fatal(err)
			}

			fields = append(fields, []byte(field))

			checkErrors(in(key.Hex(), field, txn, value, duration))
			dst := cluster.Members(bs.Key(key.Hex()))

			for e := range dst {
				if err := c.ErrorFromElement(e); err != nil {
					typex.Fatal(err)
				}
				for _, v := range c.KeysFromElement(e) {
					if parts := strings.Split(v.String(), "_"); parts[0] == field {
						return true
					}
				}
			}
			return false
		}
	)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

func TestScore(t *testing.T) {
	var (
		amount = rand.Intn(5) + 1

		creator p.RedisCreator
	)
	if defaultUseStubs {
		t.Skip("Unable to test because of concurency.")
	}

	var (
		cluster = newCluster(creator)
		in      = insert(cluster, amount)
		pool    = getIdentPool()

		f = func(txn, value string, duration time.Duration) bool {
			key, err := b.Bson(pool.Get())
			if err != nil {
				typex.Fatal(err)
			}
			field, err := b.Bson(pool.Get())
			if err != nil {
				typex.Fatal(err)
			}

			checkErrors(in(key.Hex(), field.String(), txn, value, time.Hour))

			node := selectors.KeyFieldTxnValue{
				Key:   bs.Key(key.Hex()),
				Field: bs.Key(fmt.Sprintf("%s_%d", field.String(), 0)),
				Txn:   bs.Key(txn),
				Value: value,
			}
			result, err := cluster.Score([]selectors.KeyFieldTxnValue{node})

			if err != nil {
				typex.Fatal(err)
			}
			if len(result) == 0 {
				return false
			}

			return result[node].Present
		}
	)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}
