package plaintext

import (
	"fmt"
	"io"
	"time"

	"github.com/SimonRichardson/colossus/alertmanager"
)

type manager struct {
	alertmanager.AlertBase
}

func New(w io.Writer) alertmanager.AlertManager {
	return manager{alertmanager.Make(func(s string) {
		fmt.Fprintf(w, s)
	})}
}

func (m manager) TopologyPanic() alertmanager.Cancellable {
	return m.Dispatch("topology_panic.count 1", alertmanager.AlertOptions{
		Repititions: 1,
		Delay:       time.Second * 30,
	})
}

func (m manager) CoordinatorPanic() alertmanager.Cancellable {
	return m.Dispatch("coordinator_panic.count 1", alertmanager.AlertOptions{
		Repititions: 1,
		Delay:       time.Second * 30,
	})
}
