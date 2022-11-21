package noop

import "github.com/SimonRichardson/colossus/alertmanager"

type manager struct{}

func New() alertmanager.AlertManager {
	return manager{}
}

func (m manager) TopologyPanic() alertmanager.Cancellable { return cancellable{} }

func (m manager) CoordinatorPanic() alertmanager.Cancellable { return cancellable{} }

type cancellable struct{}

func (c cancellable) Cancel() {}
