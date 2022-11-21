package selectors

import (
	"net/http"

	fs "github.com/SimonRichardson/colossus/fusion/selectors"
)

// Semaphore describes how to lock and release a operation that needs to be
// synchronised between distributed applications.
type Semaphore interface {
	Lock(Namespace) (SemaphoreUnlock, error)
}

// Heartbeat describes how to notify the reset of the datacenter that the
// application is still alive.
type Heartbeat interface {
	Heartbeat(HealthStatus) error
}

type KeyStore interface {
	List(fs.Prefix) (map[string]int, error)
}

type Request interface {
	Request(*http.Request) (*http.Response, error)
}

// Enqueuer describes how to enqueue and dequeue bytes into the system
type Enqueuer interface {
	EnqueueBytes(Queue, Class, []byte) error
	DequeueBytes(Queue, Class) ([]byte, error)
}

type Register interface {
	RegisterFailure(Queue, Class, Failure) error
}

type EventSelector interface {
	SelectEventByKey(Key) (Event, error)
	SelectEventsByOffset(int, int) ([]Event, error)
}

// Charger defines an interface for chargine items via a interface
type Charger interface {
	Charge(Event, User, Payment) (Key, error)
}

type CodeSetSelector interface {
	SelectCodeForEvent(Event, User) (CodeSet, error)
}

// Inspector provides an interface to help query the different services and
// sub-systems within bombe
type Inspector interface {
	Version() (map[string][]Version, error)
}

// Encoder describes how to encode and decode events from with in a system.
type Encoder interface {
	GetBytes(Key) ([]byte, error)
	SetBytes(Key, []byte) error
	DelBytes(Key) error
}
