package request

import "time"

type Instrumentation interface {
	RequestInstrumentation
}

type RequestInstrumentation interface {
	RequestCall()
	RequestSendTo(int)
	RequestDuration(time.Duration)
	RequestRetrieved(int)
	RequestReturned(int)
}
