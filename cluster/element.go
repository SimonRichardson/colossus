package cluster

import (
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/errors"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// ErrPartialInsertions defines an error where insertions happened, but no redis
// error replied, yet the number of actual insertions didn't match the number of
// required insertions.
var (
	ErrPartialInsertions    = typex.Errorf(errors.Source, errors.Partial, "Partial Insertions")
	ErrPartialModifications = typex.Errorf(errors.Source, errors.Partial, "Partial Modifications")
	ErrPartialDeletions     = typex.Errorf(errors.Source, errors.Partial, "Partial Deletions")
)

// ElementType defines the type of element to expect over the wire.
type ElementType int

// KeyFieldScoreTxnValueType and the following defines all the types you can expect
// when interacting with a Element.
const (
	ErrorElementType ElementType = iota
	KeyElementType
	KeyFieldSizeElementType
	KeyFieldScoreTxnValueType
	KeyFieldScoreSizeExpiryType
	CountElementType
)

// Element combines a submitted key with the resulting values. If there was an
// error while selecting a key, the error field will be populated.
type Element interface {
	Key() bs.Key
	Type() ElementType
}

// ErrorElement defines a struct that is a container for errors.
type ErrorElement struct {
	key bs.Key
	typ ElementType
	err error
}

// NewErrorElement creates a new ErrorElement
func NewErrorElement(key bs.Key, err error) *ErrorElement {
	return &ErrorElement{key, ErrorElementType, err}
}

// Key defines the key associated with the ErrorElement
func (e *ErrorElement) Key() bs.Key { return e.key }

// Type defines the type associated with the ErrorElement
func (e *ErrorElement) Type() ElementType { return e.typ }

// Error defines the error associated with the ErrorElement
func (e *ErrorElement) Error() error { return e.err }

type errorElement interface {
	Error() error
}

// ErrorFromElement attempts to get an error from the element if it exists.
func ErrorFromElement(e Element) error {
	if ee, ok := e.(errorElement); ok {
		return ee.Error()
	}
	return nil
}

// KeyFieldScoreTxnValue defines a struct that is a container for members with in the
// store.
type KeyFieldScoreTxnValue struct {
	key     bs.Key
	typ     ElementType
	members []s.KeyFieldScoreTxnValue
}

// NewKeyFieldScoreTxnValue creates a new KeyFieldScoreTxnValue
func NewKeyFieldScoreTxnValue(key bs.Key,
	members []s.KeyFieldScoreTxnValue,
) *KeyFieldScoreTxnValue {
	return &KeyFieldScoreTxnValue{key, KeyFieldScoreTxnValueType, members}
}

// Key defines the key associated with the KeyFieldScoreTxnValue
func (e *KeyFieldScoreTxnValue) Key() bs.Key { return e.key }

// Type defines the type associated with the KeyFieldScoreTxnValue
func (e *KeyFieldScoreTxnValue) Type() ElementType { return e.typ }

// Members defines the members associated with the KeyFieldScoreTxnValue
func (e *KeyFieldScoreTxnValue) Members() []s.KeyFieldScoreTxnValue { return e.members }

type keyFieldScoreTxnValue interface {
	Members() []s.KeyFieldScoreTxnValue
}

// ValuesFromElement attempts to get an key score members from the element if
// it exists.
func ValuesFromElement(e Element) []s.KeyFieldScoreTxnValue {
	if ke, ok := e.(keyFieldScoreTxnValue); ok {
		return ke.Members()
	}
	return []s.KeyFieldScoreTxnValue{}
}

// CountElement defines a struct that is a container for counting items
// (possible changes or just returning an int)
type CountElement struct {
	key    bs.Key
	typ    ElementType
	amount int
}

// NewCountElement creates a new CountElement
func NewCountElement(key bs.Key, amount int) *CountElement {
	return &CountElement{key, CountElementType, amount}
}

// Key defines the key associated with the CountElement
func (e *CountElement) Key() bs.Key { return e.key }

// Type defines the type associated with the CountElement
func (e *CountElement) Type() ElementType { return e.typ }

// Amount defines the amount associated with the CountElement
func (e *CountElement) Amount() int { return e.amount }

type amountElement interface {
	Amount() int
}

// AmountFromElement attempts to get an amount from the element if it exists.
func AmountFromElement(e Element) int {
	if ae, ok := e.(amountElement); ok {
		return ae.Amount()
	}
	return 0
}

// KeyElement defines a struct that is a container for members with in the
// store.
type KeyElement struct {
	key     bs.Key
	typ     ElementType
	members []bs.Key
}

// NewKeyElement creates a new KeyElement
func NewKeyElement(key bs.Key, keys []bs.Key) *KeyElement {
	return &KeyElement{key, KeyElementType, keys}
}

// Key defines the key associated with the KeyElement
func (e *KeyElement) Key() bs.Key { return e.key }

// Type defines the type associated with the KeyElement
func (e *KeyElement) Type() ElementType { return e.typ }

// Keys defines the key associated with the KeyElement
func (e *KeyElement) Keys() []bs.Key { return e.members }

type keyElement interface {
	Keys() []bs.Key
}

// KeysFromElement attempts to get an key score members from the element if
// it exists.
func KeysFromElement(e Element) []bs.Key {
	if ke, ok := e.(keyElement); ok {
		return ke.Keys()
	}
	return []bs.Key{}
}

// KeyFieldScoreSizeExpiryElement defines a struct that is a container for members with in the
// store.
type KeyFieldScoreSizeExpiryElement struct {
	key    bs.Key
	typ    ElementType
	member s.KeyFieldScoreSizeExpiry
}

// NewKeyFieldScoreSizeExpiryElement creates a new KeyFieldScoreSizeExpiryElement
func NewKeyFieldScoreSizeExpiryElement(keyFieldScoreSizeExpiry s.KeyFieldScoreSizeExpiry) *KeyFieldScoreSizeExpiryElement {
	return &KeyFieldScoreSizeExpiryElement{
		keyFieldScoreSizeExpiry.Key,
		KeyFieldScoreSizeExpiryType,
		keyFieldScoreSizeExpiry,
	}
}

// Key defines the key associated with the KeyFieldScoreSizeExpiryElement
func (e *KeyFieldScoreSizeExpiryElement) Key() bs.Key { return e.key }

// Type defines the type associated with the KeyFieldScoreSizeExpiryElement
func (e *KeyFieldScoreSizeExpiryElement) Type() ElementType { return e.typ }

// KeyFieldScoreSizeExpiry defines the type associated with the
// KeyFieldScoreSizeExpiryElement
func (e *KeyFieldScoreSizeExpiryElement) KeyFieldScoreSizeExpiry() s.KeyFieldScoreSizeExpiry {
	return e.member
}

type keyFieldScoreSizeExpiryElement interface {
	KeyFieldScoreSizeExpiry() s.KeyFieldScoreSizeExpiry
}

// KeyFieldScoreSizeExpiryFromElement attempts to get an key score members from
// the element if it exists.
func KeyFieldScoreSizeExpiryFromElement(e Element) s.KeyFieldScoreSizeExpiry {
	if ke, ok := e.(keyFieldScoreSizeExpiryElement); ok {
		return ke.KeyFieldScoreSizeExpiry()
	}
	return s.KeyFieldScoreSizeExpiry{}
}
