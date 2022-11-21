package consul

import "time"

type Instrumentation interface {
	SemaphoreInstrumentation
	HeartbeatInstrumentation
	KeyStoreInstrumentation
}

type SemaphoreInstrumentation interface {
	SemaphoreCall()
	SemaphoreSendTo(int)
	SemaphoreDuration(time.Duration)
	SemaphoreRetrieved(int)
	SemaphoreReturned(int)
}

type HeartbeatInstrumentation interface {
	HeartbeatCall()
	HeartbeatSendTo(int)
	HeartbeatDuration(time.Duration)
	HeartbeatRetrieved(int)
	HeartbeatReturned(int)
}

type KeyStoreInstrumentation interface {
	KeyStoreCall()
	KeyStoreSendTo(int)
	KeyStoreDuration(time.Duration)
	KeyStoreRetrieved(int)
	KeyStoreReturned(int)
}
