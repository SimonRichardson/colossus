package selectors

import (
	"strings"
	"time"

	s "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
)

// Channel defines an alias for a publishing channel
type Channel string

func (c Channel) String() string {
	return string(c)
}

// Transformer transforms an KeyFieldScoreTxnValue into a map for storing
type Transformer func(KeyFieldScoreTxnValue) (map[string]interface{}, error)

// Accessor transforms a interface{} (record) in place for modifying.
type Accessor interface {
	GetFieldValue(interface{}, string) (string, error)
	SetFieldValue(interface{}, string, string) error
}

type KeysSort []s.Key

func (a KeysSort) Len() int           { return len(a) }
func (a KeysSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a KeysSort) Less(i, j int) bool { return a[i].String() < a[j].String() }

// KeyValue pairs a key and a Value
type KeyValue struct {
	Key   s.Key
	Value string
}

// KeyField pairs a key and a field
type KeyField struct {
	Key, Field s.Key
}

// KeyFieldScoreSizeExpiry pairs a key, field, score and a size with expiry
type KeyFieldScoreSizeExpiry struct {
	Key, Field s.Key
	Score      float64
	Size       int64
	Expiry     time.Duration
}

// KeyFieldScoreSizeExpiries represents an alias for a slice of KeyFieldScoreSizeExpiry
type KeyFieldScoreSizeExpiries []KeyFieldScoreSizeExpiry

// Bucketize removes the duplicate keys so we can efficently call the storage.
func (k KeyFieldScoreSizeExpiries) Bucketize() ([]s.Key, map[s.Key][]KeyFieldScoreSizeExpiry) {
	a := map[s.Key][]KeyFieldScoreSizeExpiry{}

	for _, v := range k {
		a[v.Key] = append(a[v.Key], v)
	}

	b := make([]s.Key, 0, len(a))
	for k := range a {
		b = append(b, k)
	}

	return b, a
}

// KeyFieldSizeExpiry pairs a key, field and a size with expiry
type KeyFieldSizeExpiry struct {
	Key, Field s.Key
	Size       int64
	Expiry     time.Duration
}

// KeyFieldSizeExpiries represents an alias for a slice of KeyFieldSizeExpiry
type KeyFieldSizeExpiries []KeyFieldSizeExpiry

// Bucketize removes the duplicate keys so we can efficently call the storage.
func (k KeyFieldSizeExpiries) Bucketize() ([]s.Key, map[s.Key][]KeyFieldSizeExpiry) {
	a := map[s.Key][]KeyFieldSizeExpiry{}

	for _, v := range k {
		a[v.Key] = append(a[v.Key], v)
	}

	b := make([]s.Key, 0, len(a))
	for k := range a {
		b = append(b, k)
	}

	return b, a
}

// KeyFieldTxnValue trios a key, field, transaction and a value
type KeyFieldTxnValue struct {
	Key, Field, Txn s.Key
	Value           string
}

// KeyFieldTxnValues represents an alias for a slice of KeyFieldTxnValue
type KeyFieldTxnValues []KeyFieldTxnValue

// KeyFieldScoreTxnValues returns a KeyFieldScoreTxnValues from a KeyFieldTxnValues
func (f KeyFieldTxnValues) KeyFieldScoreTxnValues(score float64) KeyFieldScoreTxnValues {
	result := make([]KeyFieldScoreTxnValue, 0, len(f))
	for _, m := range f {
		result = append(result, KeyFieldScoreTxnValue{
			Key:   m.Key,
			Field: m.Field,
			Score: score,
			Txn:   m.Txn,
			Value: m.Value,
		})
	}
	return result
}

// FieldTxnValue pairs a field, transaction and a value
type FieldTxnValue struct {
	Field, Txn s.Key
	Value      string
}

// FieldTxnValues represents an alias for a slice of FieldTxnValue
type FieldTxnValues []FieldTxnValue

// KeyFieldScoreTxnValues returns a KeyFieldScoreTxnValue from a field and value
func (f FieldTxnValues) KeyFieldScoreTxnValues(key s.Key, score float64) KeyFieldScoreTxnValues {
	result := make([]KeyFieldScoreTxnValue, 0, len(f))
	for _, m := range f {
		result = append(result, KeyFieldScoreTxnValue{
			Key:   key,
			Field: m.Field,
			Score: score,
			Txn:   m.Txn,
			Value: m.Value,
		})
	}
	return result
}

// KeyFieldScoreTxnValue pairs a key, field, score, transaction and a Value
type KeyFieldScoreTxnValue struct {
	Key, Field s.Key
	Score      float64
	Txn        s.Key
	Value      string
}

// KeyValue returns a KeyValue from a KeyFieldScoreTxnValue
func (k KeyFieldScoreTxnValue) KeyValue() KeyValue {
	return KeyValue{Key: k.Key, Value: k.Value}
}

// KeyField returns a KeyField from a KeyFieldScoreTxnValue
func (k KeyFieldScoreTxnValue) KeyField() KeyField {
	return KeyField{Key: k.Key, Field: k.Field}
}

// KeyFieldTxnValue returns a KeyFieldTxnValue from a KeyFieldScoreTxnValue
func (k KeyFieldScoreTxnValue) KeyFieldTxnValue() KeyFieldTxnValue {
	return KeyFieldTxnValue{
		Key:   k.Key,
		Field: k.Field,
		Txn:   k.Txn,
		Value: k.Value,
	}
}

// FieldValue returns a KeyValue from a KeyFieldScoreTxnValue
func (k KeyFieldScoreTxnValue) FieldValue() KeyValue {
	return KeyValue{Key: k.Field, Value: k.Value}
}

// KeyFieldScoreTxnValues represents an alias for a slice of
// KeyFieldScoreTxnValue
type KeyFieldScoreTxnValues []KeyFieldScoreTxnValue

// Bucketize removes the duplicate keys so we can efficently call the storage.
func (k KeyFieldScoreTxnValues) Bucketize() map[s.Key][]KeyFieldScoreTxnValue {
	m := map[s.Key][]KeyFieldScoreTxnValue{}

	for _, v := range k {
		m[v.Key] = append(m[v.Key], v)
	}

	return m
}

// KeysBucketize is similar to Bucketize, except it also returns the keys as
// well
func (k KeyFieldScoreTxnValues) KeysBucketize() ([]s.Key, map[s.Key][]KeyFieldScoreTxnValue) {
	// we can state that all keys are in buckets now.
	var (
		a = k.Bucketize()
		b = make([]s.Key, 0, len(a))
	)
	for k := range a {
		b = append(b, k)
	}

	return b, a
}

// KeyValues returns a slice of KeyValue
func (k KeyFieldScoreTxnValues) KeyValues() []KeyValue {
	result := make([]KeyValue, 0, len(k))
	for _, m := range k {
		result = append(result, m.KeyValue())
	}
	return result
}

// KeyFields returns a slice of KeyField
func (k KeyFieldScoreTxnValues) KeyFields() []KeyField {
	result := make([]KeyField, 0, len(k))
	for _, m := range k {
		result = append(result, m.KeyField())
	}
	return result
}

// KeyFieldTxnValues returns a slice of KeyFieldTxnValue
func (k KeyFieldScoreTxnValues) KeyFieldTxnValues() []KeyFieldTxnValue {
	result := make([]KeyFieldTxnValue, 0, len(k))
	for _, m := range k {
		result = append(result, m.KeyFieldTxnValue())
	}
	return result
}

// FieldValues returns a slice of KeyValue
func (k KeyFieldScoreTxnValues) FieldValues() []KeyValue {
	result := make([]KeyValue, 0, len(k))
	for _, m := range k {
		result = append(result, m.FieldValue())
	}
	return result
}

// KeyFieldScoreSizeExpiry returns a slice of KeyFieldScoreSizeExpiry
func (k KeyFieldScoreTxnValues) KeyFieldScoreSizeExpiry(size KeySizeExpiry) []KeyFieldScoreSizeExpiry {
	result := make([]KeyFieldScoreSizeExpiry, 0, len(k))
	for _, v := range k {
		sizeExpiry := size[v.Key]
		result = append(result, KeyFieldScoreSizeExpiry{
			Key:    v.Key,
			Field:  v.Field,
			Score:  v.Score,
			Size:   sizeExpiry.Size,
			Expiry: sizeExpiry.Expiry,
		})
	}
	return result
}

// KeyFieldSizeExpiry returns a slice of KeyFieldSizeExpiry
func (k KeyFieldScoreTxnValues) KeyFieldSizeExpiry(size KeySizeExpiry) []KeyFieldSizeExpiry {
	result := make([]KeyFieldSizeExpiry, 0, len(k))
	for _, v := range k {
		sizeExpiry := size[v.Key]
		result = append(result, KeyFieldSizeExpiry{
			Key:    v.Key,
			Field:  v.Field,
			Size:   sizeExpiry.Size,
			Expiry: sizeExpiry.Expiry,
		})
	}
	return result
}

// KeyCount pairs a key, count
type KeyCount struct {
	Key   s.Key
	Count int
}

// Presence represents the state of a given key-Value in a cluster.
type Presence struct {
	Present  bool
	Inserted bool
	Score    float64
}

type SizeExpiry struct {
	Size   int64
	Expiry time.Duration
}

// KeySizeExpiry represents a pair of Keys and Sizes with Expiry time
type KeySizeExpiry map[s.Key]SizeExpiry

// MakeKeySizeExpiry creates a new KeySizeExpiry
func MakeKeySizeExpiry() KeySizeExpiry {
	return map[s.Key]SizeExpiry{}
}

// Get returns a possible size or an error if a size associated with a key isn't
// found.
func (k KeySizeExpiry) Get(key s.Key) (SizeExpiry, error) {
	if v, ok := k[key]; ok {
		return v, nil
	}
	return SizeExpiry{}, typex.Errorf(errors.Source, errors.NoCaseFound, "Not found")
}

// MakeKeySizeSingleton creates a KeySizeExpiry with one element.
func MakeKeySizeSingleton(key s.Key, maxSize int64, expiry time.Duration) KeySizeExpiry {
	return map[s.Key]SizeExpiry{
		key: SizeExpiry{maxSize, expiry},
	}
}

// Operation describes a how to manage patches to the stores, with expectations.
type Operation struct {
	Op    Op
	Path  Path
	Value string
}

// Op defines a typed alias for all the operations
type Op string

func (o Op) String() string {
	return string(o)
}

// Path defines a type alias for describing locations to modify
type Path string

// Parts returns the indvidual parts of the path
func (p Path) Parts() []string {
	return strings.Split(p.String(), "/")
}

func (p Path) String() string {
	return string(p)
}

// QueryOptions defines items that we can query on.
type QueryOptions struct {
	OwnerId s.Key
}

type QueryRecord struct {
	Key    s.Key
	Field  s.Key
	Record map[string]interface{}
}
