package selectors

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"testing/quick"
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
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

func config() *quick.Config {
	if testing.Short() {
		return &quick.Config{
			MaxCount:      10,
			MaxCountScale: 10,
		}
	}
	return nil
}

func TestKey_Len(t *testing.T) {
	var (
		f = func(s string) int {
			return bs.Key(s).Len()
		}
		g = func(s string) int {
			return len(s)
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKey_String(t *testing.T) {
	var (
		f = func(s string) string {
			return bs.Key(s).String()
		}
		g = func(s string) string {
			return s
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

// KeyFieldSizeExpiry

func TestKeyFieldScoreSizeExpiries_Bucketize(t *testing.T) {
	var (
		f = func(key, field string, score float64, size int64, expiry time.Duration) ([]bs.Key, map[bs.Key][]KeyFieldScoreSizeExpiry) {
			item := KeyFieldScoreSizeExpiry{
				Key:    bs.Key(key),
				Field:  bs.Key(field),
				Score:  score,
				Size:   size,
				Expiry: expiry,
			}
			items := KeyFieldScoreSizeExpiries([]KeyFieldScoreSizeExpiry{
				item,
			})
			return items.Bucketize()
		}
		g = func(key, field string, score float64, size int64, expiry time.Duration) ([]bs.Key, map[bs.Key][]KeyFieldScoreSizeExpiry) {
			k := bs.Key(key)
			return []bs.Key{k}, map[bs.Key][]KeyFieldScoreSizeExpiry{
				k: []KeyFieldScoreSizeExpiry{
					KeyFieldScoreSizeExpiry{
						Key:    k,
						Field:  bs.Key(field),
						Score:  score,
						Size:   size,
						Expiry: expiry,
					},
				},
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreSizeExpiries_BucketizeMultipleItems(t *testing.T) {
	var (
		f = func(keys []string, field string, score float64, size int64, expiry time.Duration) ([]bs.Key, map[bs.Key][]KeyFieldScoreSizeExpiry) {
			values := make([]KeyFieldScoreSizeExpiry, 0, len(keys))
			for _, v := range keys {
				values = append(values, KeyFieldScoreSizeExpiry{
					Key:    bs.Key(v),
					Field:  bs.Key(field),
					Score:  score,
					Size:   size,
					Expiry: expiry,
				})
			}
			items := KeyFieldScoreSizeExpiries(values)
			a, b := items.Bucketize()
			sort.Sort(KeysSort(a))
			return a, b
		}
		g = func(keys []string, field string, score float64, size int64, expiry time.Duration) ([]bs.Key, map[bs.Key][]KeyFieldScoreSizeExpiry) {
			b := map[bs.Key][]KeyFieldScoreSizeExpiry{}
			for _, v := range keys {
				var (
					key  = bs.Key(v)
					item = KeyFieldScoreSizeExpiry{
						Key:    key,
						Field:  bs.Key(field),
						Score:  score,
						Size:   size,
						Expiry: expiry,
					}
				)
				b[item.Key] = append(b[item.Key], item)
			}

			a := make([]bs.Key, 0, len(b))
			for k := range b {
				a = append(a, k)
			}
			sort.Sort(KeysSort(a))
			return a, b
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

// KeyFieldSizeExpiry

func TestKeyFieldSizeExpiries_Bucketize(t *testing.T) {
	var (
		f = func(key, field string, size int64, expiry time.Duration) ([]bs.Key, map[bs.Key][]KeyFieldSizeExpiry) {
			item := KeyFieldSizeExpiry{
				Key:    bs.Key(key),
				Field:  bs.Key(field),
				Size:   size,
				Expiry: expiry,
			}
			items := KeyFieldSizeExpiries([]KeyFieldSizeExpiry{
				item,
			})
			return items.Bucketize()
		}
		g = func(key, field string, size int64, expiry time.Duration) ([]bs.Key, map[bs.Key][]KeyFieldSizeExpiry) {
			k := bs.Key(key)
			return []bs.Key{k}, map[bs.Key][]KeyFieldSizeExpiry{
				k: []KeyFieldSizeExpiry{
					KeyFieldSizeExpiry{
						Key:    k,
						Field:  bs.Key(field),
						Size:   size,
						Expiry: expiry,
					},
				},
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldSizeExpiries_BucketizeMultipleItems(t *testing.T) {
	var (
		f = func(keys []string, field string, size int64, expiry time.Duration) ([]bs.Key, map[bs.Key][]KeyFieldSizeExpiry) {
			values := make([]KeyFieldSizeExpiry, 0, len(keys))
			for _, v := range keys {
				values = append(values, KeyFieldSizeExpiry{
					Key:    bs.Key(v),
					Field:  bs.Key(field),
					Size:   size,
					Expiry: expiry,
				})
			}
			items := KeyFieldSizeExpiries(values)
			a, b := items.Bucketize()
			sort.Sort(KeysSort(a))
			return a, b
		}
		g = func(keys []string, field string, size int64, expiry time.Duration) ([]bs.Key, map[bs.Key][]KeyFieldSizeExpiry) {
			b := map[bs.Key][]KeyFieldSizeExpiry{}
			for _, v := range keys {
				var (
					key  = bs.Key(v)
					item = KeyFieldSizeExpiry{
						Key:    key,
						Field:  bs.Key(field),
						Size:   size,
						Expiry: expiry,
					}
				)
				b[item.Key] = append(b[item.Key], item)
			}

			a := make([]bs.Key, 0, len(b))
			for k := range b {
				a = append(a, k)
			}
			sort.Sort(KeysSort(a))
			return a, b
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

// KeyKeyFieldTxnValue

func TestKeyFieldTxnValue_KeyFieldScoreTxnValues(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) KeyFieldScoreTxnValues {
			item := KeyFieldTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Txn:   bs.Key(txn),
				Value: value,
			}
			items := KeyFieldTxnValues([]KeyFieldTxnValue{
				item,
			})
			return items.KeyFieldScoreTxnValues(score)
		}
		g = func(key, field string, score float64, txn, value string) KeyFieldScoreTxnValues {
			return KeyFieldScoreTxnValues([]KeyFieldScoreTxnValue{
				KeyFieldScoreTxnValue{
					Key:   bs.Key(key),
					Field: bs.Key(field),
					Score: score,
					Txn:   bs.Key(txn),
					Value: value,
				},
			})
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

// FieldTxnValue

func TestFieldTxnValue_KeyFieldScoreTxnValues(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) KeyFieldScoreTxnValues {
			item := FieldTxnValue{
				Field: bs.Key(field),
				Txn:   bs.Key(txn),
				Value: value,
			}
			items := FieldTxnValues([]FieldTxnValue{
				item,
			})
			return items.KeyFieldScoreTxnValues(bs.Key(key), score)
		}
		g = func(key, field string, score float64, txn, value string) KeyFieldScoreTxnValues {
			return KeyFieldScoreTxnValues([]KeyFieldScoreTxnValue{
				KeyFieldScoreTxnValue{
					Key:   bs.Key(key),
					Field: bs.Key(field),
					Score: score,
					Txn:   bs.Key(txn),
					Value: value,
				},
			})
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

// KeyFieldScoreTxnValue

func TestKeyFieldScoreTxnValue_KeyValue(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) KeyValue {
			item := KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Score: score,
				Txn:   bs.Key(txn),
				Value: value,
			}
			return item.KeyValue()
		}
		g = func(key, field string, score float64, txn, value string) KeyValue {
			return KeyValue{
				Key:   bs.Key(key),
				Value: value,
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreTxnValue_KeyField(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) KeyField {
			item := KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Score: score,
				Txn:   bs.Key(txn),
				Value: value,
			}
			return item.KeyField()
		}
		g = func(key, field string, score float64, txn, value string) KeyField {
			return KeyField{
				Key:   bs.Key(key),
				Field: bs.Key(field),
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreTxnValue_KeyFieldTxnValue(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) KeyFieldTxnValue {
			item := KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Score: score,
				Txn:   bs.Key(txn),
				Value: value,
			}
			return item.KeyFieldTxnValue()
		}
		g = func(key, field string, score float64, txn, value string) KeyFieldTxnValue {
			return KeyFieldTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Txn:   bs.Key(txn),
				Value: value,
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreTxnValue_FieldValue(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) KeyValue {
			item := KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Score: score,
				Txn:   bs.Key(txn),
				Value: value,
			}
			return item.FieldValue()
		}
		g = func(key, field string, score float64, txn, value string) KeyValue {
			return KeyValue{
				Key:   bs.Key(field),
				Value: value,
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreValue_Bucketize(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) map[bs.Key][]KeyFieldScoreTxnValue {
			item := KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Score: score,
				Txn:   bs.Key(txn),
				Value: value,
			}
			items := KeyFieldScoreTxnValues([]KeyFieldScoreTxnValue{
				item,
			})
			return items.Bucketize()
		}
		g = func(key, field string, score float64, txn, value string) map[bs.Key][]KeyFieldScoreTxnValue {
			k := bs.Key(key)
			return map[bs.Key][]KeyFieldScoreTxnValue{
				k: []KeyFieldScoreTxnValue{
					KeyFieldScoreTxnValue{
						Key:   k,
						Field: bs.Key(field),
						Score: score,
						Txn:   bs.Key(txn),
						Value: value,
					},
				},
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreValue_KeysBucketize(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) ([]bs.Key, map[bs.Key][]KeyFieldScoreTxnValue) {
			item := KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Score: score,
				Txn:   bs.Key(txn),
				Value: value,
			}
			items := KeyFieldScoreTxnValues([]KeyFieldScoreTxnValue{
				item,
			})
			return items.KeysBucketize()
		}
		g = func(key, field string, score float64, txn, value string) ([]bs.Key, map[bs.Key][]KeyFieldScoreTxnValue) {
			k := bs.Key(key)
			return []bs.Key{k}, map[bs.Key][]KeyFieldScoreTxnValue{
				k: []KeyFieldScoreTxnValue{
					KeyFieldScoreTxnValue{
						Key:   k,
						Field: bs.Key(field),
						Score: score,
						Txn:   bs.Key(txn),
						Value: value,
					},
				},
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreValue_KeysBucketizeMultipleItems(t *testing.T) {
	var (
		f = func(keys []string, field string, score float64, txn, value string) ([]bs.Key, map[bs.Key][]KeyFieldScoreTxnValue) {
			values := make([]KeyFieldScoreTxnValue, 0, len(keys))
			for _, v := range keys {
				values = append(values, KeyFieldScoreTxnValue{
					Key:   bs.Key(v),
					Field: bs.Key(field),
					Score: score,
					Txn:   bs.Key(txn),
					Value: value,
				})
			}
			items := KeyFieldScoreTxnValues(values)
			a, b := items.KeysBucketize()
			sort.Sort(KeysSort(a))
			return a, b
		}
		g = func(keys []string, field string, score float64, txn, value string) ([]bs.Key, map[bs.Key][]KeyFieldScoreTxnValue) {
			b := map[bs.Key][]KeyFieldScoreTxnValue{}
			for _, v := range keys {
				var (
					key  = bs.Key(v)
					item = KeyFieldScoreTxnValue{
						Key:   key,
						Field: bs.Key(field),
						Score: score,
						Txn:   bs.Key(txn),
						Value: value,
					}
				)
				b[item.Key] = append(b[item.Key], item)
			}

			a := make([]bs.Key, 0, len(b))
			for k := range b {
				a = append(a, k)
			}
			sort.Sort(KeysSort(a))
			return a, b
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreTxnValue_KeyValues(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) []KeyValue {
			item := KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Score: score,
				Txn:   bs.Key(txn),
				Value: value,
			}
			items := KeyFieldScoreTxnValues([]KeyFieldScoreTxnValue{
				item,
			})
			return items.KeyValues()
		}
		g = func(key, field string, score float64, txn, value string) []KeyValue {
			return []KeyValue{
				KeyValue{
					Key:   bs.Key(key),
					Value: value,
				},
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreTxnValue_KeyFields(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) []KeyField {
			item := KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Score: score,
				Txn:   bs.Key(txn),
				Value: value,
			}
			items := KeyFieldScoreTxnValues([]KeyFieldScoreTxnValue{
				item,
			})
			return items.KeyFields()
		}
		g = func(key, field string, score float64, txn, value string) []KeyField {
			return []KeyField{
				KeyField{
					Key:   bs.Key(key),
					Field: bs.Key(field),
				},
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreTxnValue_KeyFieldTxnValues(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) []KeyFieldTxnValue {
			item := KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Score: score,
				Txn:   bs.Key(txn),
				Value: value,
			}
			items := KeyFieldScoreTxnValues([]KeyFieldScoreTxnValue{
				item,
			})
			return items.KeyFieldTxnValues()
		}
		g = func(key, field string, score float64, txn, value string) []KeyFieldTxnValue {
			return []KeyFieldTxnValue{
				KeyFieldTxnValue{
					Key:   bs.Key(key),
					Field: bs.Key(field),
					Txn:   bs.Key(txn),
					Value: value,
				},
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

func TestKeyFieldScoreTxnValue_FieldValues(t *testing.T) {
	var (
		f = func(key, field string, score float64, txn, value string) []KeyValue {
			item := KeyFieldScoreTxnValue{
				Key:   bs.Key(key),
				Field: bs.Key(field),
				Score: score,
				Txn:   bs.Key(txn),
				Value: value,
			}
			items := KeyFieldScoreTxnValues([]KeyFieldScoreTxnValue{
				item,
			})
			return items.FieldValues()
		}
		g = func(key, field string, score float64, txn, value string) []KeyValue {
			return []KeyValue{
				KeyValue{
					Key:   bs.Key(field),
					Value: value,
				},
			}
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}

// Path

func TestPath_Parts(t *testing.T) {
	var (
		f = func(a, b string) []string {
			return Path(fmt.Sprintf("%s/%s", a, b)).Parts()
		}
		g = func(a, b string) []string {
			parts := []string{}
			parts = append(parts, strings.Split(a, "/")...)
			parts = append(parts, strings.Split(b, "/")...)
			return parts
		}
	)

	if err := quick.CheckEqual(f, g, config()); err != nil {
		t.Error(err)
	}
}
