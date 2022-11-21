package resque

import "time"

type Instrumentation interface {
	EnqueueInstrumentation
	DequeuerInstrumentation
	RegisterInstrumentation
}

type EnqueueInstrumentation interface {
	EnqueueCall()
	EnqueueSendTo(int)
	EnqueueDuration(time.Duration)
}

type DequeuerInstrumentation interface {
	DequeueCall()
	DequeueSendTo(int)
	DequeueDuration(time.Duration)
}

type RegisterInstrumentation interface {
	RegisterFailureCall()
	RegisterFailureSendTo(int)
	RegisterFailureDuration(time.Duration)
}
