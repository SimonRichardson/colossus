package services

import (
	"net/http"

	"github.com/SimonRichardson/colossus/blist/selectors"
	fs "github.com/SimonRichardson/colossus/fusion/selectors"
)

// Semaphore represents a way to gain a distributed lock and release it when
// done, so that we can make sure that only one thing is ever done at once.
type Semaphore interface {
	Lock(selectors.Namespace) <-chan Element
}

// Heartbeat describes how to notify the reset of the datacenter that the
// application is still alive.
type Heartbeat interface {
	Heartbeat(selectors.HealthStatus) <-chan Element
}

// KeyStore describes how to notify the reset of the datacenter that the
// application is still alive.
type KeyStore interface {
	List(fs.Prefix) <-chan Element
}

type Request interface {
	Request(*http.Request) <-chan Element
}

type Index interface {
	Index() int
}

// Enqueuer represents a way to enqueue and dequeue plain items on to the
// message bus.
type Enqueuer interface {
	EnqueueBytes(selectors.Queue, selectors.Class, []byte) <-chan Element
	DequeueBytes(selectors.Queue, selectors.Class) <-chan Element
}

type Register interface {
	RegisterFailure(selectors.Queue, selectors.Class, selectors.Failure) <-chan Element
}

// Inspector represents a way check various attributes of a service.
type Inspector interface {
	Version() <-chan Element
}

// EventSelector represents a way to select events or event attributes from the
// backing store.
type EventSelector interface {
	SelectEventByKey(selectors.Key) <-chan Element
	SelectEventsByOffset(int, int) <-chan Element
}

// CodeSelector represents a way to select events or event attributes from the
// backing store.
type CodeSelector interface {
	SelectCodeForEvent(selectors.Event, selectors.User) <-chan Element
}

// Charger represents a way to purchase items from the service.
type Charger interface {
	Charge(selectors.Event, selectors.User, selectors.Payment) <-chan Element
}

// Encoder represents a way to enqueue items into a message bus.
type Encoder interface {
	GetBytes(selectors.Key) <-chan Element
	SetBytes(selectors.Key, []byte) <-chan Element
	DelBytes(selectors.Key) <-chan Element
}
