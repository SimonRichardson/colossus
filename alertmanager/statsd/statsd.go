package statsd

import (
	"time"

	"github.com/SimonRichardson/colossus/alertmanager"
	"github.com/peterbourgon/g2s"
)

type manager struct {
	alertmanager.AlertBase
}

func New(statter g2s.Statter, sampleRate float32) alertmanager.AlertManager {
	return manager{alertmanager.Make(func(s string) {
		statter.Counter(sampleRate, s, 1)
	})}
}

func (m manager) TopologyPanic() alertmanager.Cancellable {
	return m.Dispatch("topology_panic.count", alertmanager.AlertOptions{
		Repititions: 1,
		Delay:       time.Second * 30,
	})
}

func (m manager) CoordinatorPanic() alertmanager.Cancellable {
	return m.Dispatch("coordinator_panic.count", alertmanager.AlertOptions{
		Repititions: 1,
		Delay:       time.Second * 30,
	})
}
