package alertmanager

import "time"

type AlertManager interface {
	TopologyAlertManager
	CoordinatorAlertManager
}

type Cancellable interface {
	Cancel()
}

type AlertOptions struct {
	Repititions uint
	Delay       time.Duration
}

type TopologyAlertManager interface {
	TopologyPanic() Cancellable
}

type CoordinatorAlertManager interface {
	CoordinatorPanic() Cancellable
}
