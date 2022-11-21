package cache

import "time"

type Instrumentation interface {
	EncoderInstrumentation
}

type EncoderInstrumentation interface {
	EncodeCall()
	EncodeSendTo(int)
	EncodeDuration(time.Duration)
}
