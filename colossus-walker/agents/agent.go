package agents

import "github.com/SimonRichardson/colossus/coordinator"

// Agent defines a interface for creating agents that are *currently*
// unsupervised.
type Agent interface {
	Init(AgentOptions) error
}

type AgentOptions struct {
	Coordinator *coordinator.Coordinator
	HttpAddress string
}
