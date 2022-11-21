package multi

import "github.com/SimonRichardson/colossus/alertmanager"

type manager struct {
	managers []alertmanager.AlertManager
}

func New(managers ...alertmanager.AlertManager) alertmanager.AlertManager {
	return manager{managers}
}

func (m manager) TopologyPanic() alertmanager.Cancellable {
	nodes := make([]alertmanager.Cancellable, 0, len(m.managers))
	for _, v := range m.managers {
		nodes = append(nodes, v.TopologyPanic())
	}
	return &cancellable{nodes}
}

func (m manager) CoordinatorPanic() alertmanager.Cancellable {
	nodes := make([]alertmanager.Cancellable, 0, len(m.managers))
	for _, v := range m.managers {
		nodes = append(nodes, v.CoordinatorPanic())
	}
	return &cancellable{nodes}
}

type cancellable struct {
	nodes []alertmanager.Cancellable
}

func (c *cancellable) Cancel() {
	for _, v := range c.nodes {
		v.Cancel()
	}
	c.nodes = make([]alertmanager.Cancellable, 0)
}
