package consul

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/selectors"
	fs "github.com/SimonRichardson/colossus/fusion/selectors"
)

// NoopSemaphore defines a selector that performs no operations, but attemps to
// provide "sane" data that will allow the application to still execute.
func NoopSemaphore(s *Service, t Tactic) selectors.Semaphore { return noop{s} }

// NoopHeartbeat defines a selector that performs no operations, but attemps to
// provide "sane" data that will allow the application to still execute.
func NoopHeartbeat(s *Service, t Tactic) selectors.Heartbeat { return noop{s} }

// NoopKeyStore defines a selector that performs no operations, but attemps to
// provide "sane" data that will allow the application to still execute.
func NoopKeyStore(s *Service, t Tactic) selectors.KeyStore { return noop{s} }

type noop struct {
	*Service
}

func (n noop) Lock(selectors.Namespace) (selectors.SemaphoreUnlock, error) {
	return noopUnlock, nil
}

func (n noop) Heartbeat(selectors.HealthStatus) error {
	return nil
}

func (n noop) List(fs.Prefix) (map[string]int, error) {
	return nil, nil
}

func noopUnlock() error {
	return nil
}

type NoopInstrumentation struct{}

func (NoopInstrumentation) SemaphoreCall()                  {}
func (NoopInstrumentation) SemaphoreSendTo(int)             {}
func (NoopInstrumentation) SemaphoreDuration(time.Duration) {}
func (NoopInstrumentation) SemaphoreRetrieved(int)          {}
func (NoopInstrumentation) SemaphoreReturned(int)           {}
func (NoopInstrumentation) HeartbeatCall()                  {}
func (NoopInstrumentation) HeartbeatSendTo(int)             {}
func (NoopInstrumentation) HeartbeatDuration(time.Duration) {}
func (NoopInstrumentation) HeartbeatRetrieved(int)          {}
func (NoopInstrumentation) HeartbeatReturned(int)           {}
func (NoopInstrumentation) KeyStoreCall()                   {}
func (NoopInstrumentation) KeyStoreSendTo(int)              {}
func (NoopInstrumentation) KeyStoreDuration(time.Duration)  {}
func (NoopInstrumentation) KeyStoreRetrieved(int)           {}
func (NoopInstrumentation) KeyStoreReturned(int)            {}
