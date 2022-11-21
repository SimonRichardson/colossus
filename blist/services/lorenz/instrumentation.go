package lorenz

import "time"

type Instrumentation interface {
	EventSelectorInstrumentation
	CodeSelectorInstrumentation
	InspectorInstrumentation
	ChargerInstrumentation
}

type EventSelectorInstrumentation interface {
	EventSelectCall()
	EventSelectSendTo(int)
	EventSelectDuration(time.Duration)
}

type CodeSelectorInstrumentation interface {
	CodeSelectCall()
	CodeSelectSendTo(int)
	CodeSelectDuration(time.Duration)
}

type InspectorInstrumentation interface {
	InspectCall()
	InspectSendTo(int)
	InspectDuration(time.Duration)
	InspectRetrieved(int)
	InspectReturned(int)
}

type ChargerInstrumentation interface {
	ChargeCall()
	ChargeSendTo(int)
	ChargeDuration(time.Duration)
	ChargeRetrieved(int)
	ChargeReturned(int)
}
