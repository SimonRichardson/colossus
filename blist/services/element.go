package services

import (
	"net/http"

	s "github.com/SimonRichardson/colossus/blist/selectors"
)

// ElementType defines the type of element to expect over the wire.
type ElementType int

// KeyFieldScoreTxnValueType and the following defines all the types you can expect
// when interacting with a Element.
const (
	ErrorElementType ElementType = iota
	KeyElementType
	KeyTxnElementType
	ResponseElementType
	SemaphoreUnlockElementType
	MapStringIntElementType
	BytesElementType
	VersionElementType
	EventElementType
	EventsElementType
	CodeSetElementType
)

// Element combines a submitted key with the resulting values. If there was an
// error while selecting a key, the error field will be populated.
type Element interface {
	Type() ElementType
}

// ErrorElement defines a struct that is a container for errors.
type ErrorElement struct {
	typ ElementType
	err error
}

// NewErrorElement creates a new ErrorElement
func NewErrorElement(err error) *ErrorElement {
	return &ErrorElement{ErrorElementType, err}
}

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

// KeyElement defines a struct that is a container for key items
type KeyElement struct {
	key s.Key
	typ ElementType
}

// NewKeyElement creates a new KeyElement
func NewKeyElement(key s.Key) *KeyElement {
	return &KeyElement{typ: KeyElementType, key: key}
}

// Key defines the key associated with the KeyElement
func (e *KeyElement) Key() s.Key { return e.key }

// Type defines the type associated with the KeyElement
func (e *KeyElement) Type() ElementType { return e.typ }

type keyElement interface {
	Key() s.Key
}

// KeyFromElement attempts to get an keys from the element if it exists.
func KeyFromElement(e Element) s.Key {
	if ae, ok := e.(keyElement); ok {
		return ae.Key()
	}
	return s.Key("")
}

// KeyTxnElement defines a struct that is a container for key items
type KeyTxnElement struct {
	key, txn s.Key
	typ      ElementType
}

// NewKeyTxnElement creates a new KeyTxnElement
func NewKeyTxnElement(key, txn s.Key) *KeyTxnElement {
	return &KeyTxnElement{typ: KeyTxnElementType, key: key, txn: txn}
}

// Key defines the key associated with the KeyTxnElement
func (e *KeyTxnElement) Key() s.Key { return e.key }

// Type defines the type associated with the KeyTxnElement
func (e *KeyTxnElement) Type() ElementType { return e.typ }

type keyTxnElement interface {
	Txn() s.Key
}

// KeyTxnFromElement attempts to get an keys from the element if it exists.
func KeyTxnFromElement(e Element) s.Key {
	if ae, ok := e.(keyTxnElement); ok {
		return ae.Txn()
	}
	return s.Key("")
}

// ResponseElement defines a struct that is a container for response items
type ResponseElement struct {
	index    int
	response *http.Response
	typ      ElementType
}

// NewResponseElement creates a new ResponseElement
func NewResponseElement(index int, response *http.Response) *ResponseElement {
	return &ResponseElement{
		typ:      ResponseElementType,
		index:    index,
		response: response,
	}
}

// Response defines the response associated with the ResponseElement
func (e *ResponseElement) Response() *http.Response { return e.response }

// Index defines the index associated with the ResponseElement
func (e *ResponseElement) Index() int { return e.index }

// Type defines the type associated with the ResponseElement
func (e *ResponseElement) Type() ElementType { return e.typ }

type responseElement interface {
	Response() *http.Response
}

// ResponseFromElement attempts to get an responses from the element if it exists.
func ResponseFromElement(e Element) *http.Response {
	if ae, ok := e.(responseElement); ok {
		return ae.Response()
	}
	return nil
}

type indexElement interface {
	Index() int
}

// IndexFromElement attempts to get an index from the element if it exists.
func IndexFromElement(e Element) int {
	if ae, ok := e.(indexElement); ok {
		return ae.Index()
	}
	return -1
}

// SemaphoreUnlockElement defines a struct that is a container for unlock items
type SemaphoreUnlockElement struct {
	typ    ElementType
	unlock s.SemaphoreUnlock
}

// NewSemaphoreUnlockElement creates a new SemaphoreUnlockElement
func NewSemaphoreUnlockElement(unlock s.SemaphoreUnlock) *SemaphoreUnlockElement {
	return &SemaphoreUnlockElement{SemaphoreUnlockElementType, unlock}
}

// Type defines the type associated with the SemaphoreUnlockElement
func (e *SemaphoreUnlockElement) Type() ElementType { return e.typ }

// SemaphoreUnlock defines the SemaphoreUnlock associated with the SemaphoreUnlockElement
func (e *SemaphoreUnlockElement) SemaphoreUnlock() s.SemaphoreUnlock { return e.unlock }

type unlockElement interface {
	SemaphoreUnlock() s.SemaphoreUnlock
}

// SemaphoreUnlockFromElement attempts to get an SemaphoreUnlock from the element if it exists.
func SemaphoreUnlockFromElement(e Element) s.SemaphoreUnlock {
	if ae, ok := e.(unlockElement); ok {
		return ae.SemaphoreUnlock()
	}
	return func() error { return nil }
}

// MapStringIntElement defines a struct that is a container for unlock items
type MapStringIntElement struct {
	typ    ElementType
	unlock map[string]int
}

// NewMapStringIntElement creates a new MapStringIntElement
func NewMapStringIntElement(unlock map[string]int) *MapStringIntElement {
	return &MapStringIntElement{MapStringIntElementType, unlock}
}

// Type defines the type associated with the MapStringIntElement
func (e *MapStringIntElement) Type() ElementType { return e.typ }

// MapStringInt defines the MapStringInt associated with the MapStringIntElement
func (e *MapStringIntElement) MapStringInt() map[string]int { return e.unlock }

type mapStringIntElement interface {
	MapStringInt() map[string]int
}

// MapStringIntFromElement attempts to get an MapStringInt from the element if it exists.
func MapStringIntFromElement(e Element) map[string]int {
	if ae, ok := e.(mapStringIntElement); ok {
		return ae.MapStringInt()
	}
	return map[string]int{}
}

// BytesElement defines a struct that is a container for value items
type BytesElement struct {
	value []byte
	typ   ElementType
}

// NewBytesElement creates a new BytesElement
func NewBytesElement(value []byte) *BytesElement {
	return &BytesElement{typ: BytesElementType, value: value}
}

// Bytes defines the value associated with the BytesElement
func (e *BytesElement) Bytes() []byte { return e.value }

// Type defines the type associated with the BytesElement
func (e *BytesElement) Type() ElementType { return e.typ }

type valueElement interface {
	Bytes() []byte
}

// BytesFromElement attempts to get an values from the element if it exists.
func BytesFromElement(e Element) []byte {
	if ae, ok := e.(valueElement); ok {
		return ae.Bytes()
	}
	return nil
}

// VersionElement defines a struct that is a container for version items
type VersionElement struct {
	typ     ElementType
	version s.Version
}

// NewVersionElement creates a new VersionElement
func NewVersionElement(version s.Version) *VersionElement {
	return &VersionElement{VersionElementType, version}
}

// Type defines the type associated with the VersionElement
func (e *VersionElement) Type() ElementType { return e.typ }

// Version defines the version associated with the VersionElement
func (e *VersionElement) Version() s.Version { return e.version }

type versionElement interface {
	Version() s.Version
}

// VersionFromElement attempts to get an version from the element if it exists.
func VersionFromElement(e Element) s.Version {
	if ae, ok := e.(versionElement); ok {
		return ae.Version()
	}
	return s.Version("")
}

// EventElement defines a struct that is a container for event items
type EventElement struct {
	typ   ElementType
	event s.Event
}

// NewEventElement creates a new EventElement
func NewEventElement(event s.Event) *EventElement {
	return &EventElement{EventElementType, event}
}

// Type defines the type associated with the EventElement
func (e *EventElement) Type() ElementType { return e.typ }

// Event defines the Event associated with the EventElement
func (e *EventElement) Event() s.Event { return e.event }

type eventElement interface {
	Event() s.Event
}

// EventFromElement attempts to get an Event from the element if it exists.
func EventFromElement(e Element) s.Event {
	if ae, ok := e.(eventElement); ok {
		return ae.Event()
	}
	return s.Event{}
}

// EventsElement defines a struct that is a container for events items
type EventsElement struct {
	typ    ElementType
	events []s.Event
}

// NewEventsElement creates a new EventsElement
func NewEventsElement(events []s.Event) *EventsElement {
	return &EventsElement{EventsElementType, events}
}

// Type defines the type associated with the EventsElement
func (e *EventsElement) Type() ElementType { return e.typ }

// Events defines the Events associated with the EventsElement
func (e *EventsElement) Events() []s.Event { return e.events }

type eventsElement interface {
	Events() []s.Event
}

// EventsFromElement attempts to get an Events from the element if it exists.
func EventsFromElement(e Element) []s.Event {
	if ae, ok := e.(eventsElement); ok {
		return ae.Events()
	}
	return nil
}

// CodeSetElement defines a struct that is a container for CodeSet items
type CodeSetElement struct {
	typ  ElementType
	code s.CodeSet
}

// NewCodeSetElement creates a new CodeSetElement
func NewCodeSetElement(code s.CodeSet) *CodeSetElement {
	return &CodeSetElement{CodeSetElementType, code}
}

// Type defines the type associated with the CodeSetElement
func (e *CodeSetElement) Type() ElementType { return e.typ }

// CodeSet defines the CodeSet associated with the CodeSetElement
func (e *CodeSetElement) CodeSet() s.CodeSet { return e.code }

type codeSetElement interface {
	CodeSet() s.CodeSet
}

// CodeSetFromElement attempts to get an CodeSet from the element if it exists.
func CodeSetFromElement(e Element) s.CodeSet {
	if ae, ok := e.(codeSetElement); ok {
		return ae.CodeSet()
	}
	return s.CodeSet{}
}
