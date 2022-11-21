package prometheus

import (
	"time"

	"github.com/SimonRichardson/colossus/alertmanager"
	"github.com/prometheus/client_golang/prometheus"
)

type namespace string

const (
	topologyPanic    namespace = "topology_panic"
	coordinatorPanic namespace = "coordinator_panic"
)

func (n namespace) String() string {
	return string(n)
}

type manager struct {
	alertmanager.AlertBase
	counters map[string]prometheus.Counter
}

func New(prefix string, maxSummaryAge time.Duration) alertmanager.AlertManager {
	counters := map[string]prometheus.Counter{
		topologyPanic.String(): prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      topologyPanic.String(),
			Help:      "How many topology_panic have been made.",
		}),
		coordinatorPanic.String(): prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      coordinatorPanic.String(),
			Help:      "How many coordinator_panic have been made.",
		}),
	}

	for _, v := range counters {
		prometheus.MustRegister(v)
	}

	return manager{alertmanager.Make(func(s string) {
		if counter, ok := counters[s]; ok {
			counter.Inc()
		}
	}), counters}
}

func (m manager) TopologyPanic() alertmanager.Cancellable {
	return m.Dispatch(topologyPanic.String(), alertmanager.AlertOptions{
		Repititions: 3,
		Delay:       time.Second * 30,
	})
}

func (m manager) CoordinatorPanic() alertmanager.Cancellable {
	return m.Dispatch(coordinatorPanic.String(), alertmanager.AlertOptions{
		Repititions: 3,
		Delay:       time.Second * 30,
	})
}
