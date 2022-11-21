package consul

import (
	"strings"
	"time"

	"github.com/SimonRichardson/colossus/blist/common"
	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/blist/services/consul/client"
	"github.com/SimonRichardson/colossus/typex"
)

// DefaultConfig creates a default configuration for the service
func DefaultConfig(address, checkId, output string) ([]Cluster, error) {
	return ParseString(address, checkId, output, 100, nil)
}

var (
	empty = []Cluster{}
)

// StrategyOptions defines what options are available when creating a and using
// a stragegy.
type StrategyOptions struct {
	Strategy            string
	Tactic              string
	RequestsPerDuration int
	RequestsDuration    string
	Quorum              float64
}

// ParseString parses various inputs and returns a slice of clusters that we can
// then use with in the service.
func ParseString(address, checkId, output string,
	maxSize int,
	creator client.ClientCreator,
) ([]Cluster, error) {
	var (
		clusters = []Cluster{}
		fn       = clientCreator(creator)
	)

	for _, v := range strings.Split(common.StripWhitespace(address), ";") {
		clusters = append(clusters, newCluster(
			client.New(v, checkId, output, maxSize, fn),
		))
	}

	if len(clusters) < 1 {
		return empty, typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"No clusters specified %q for consul", address)
	}

	return clusters, nil
}

func clientCreator(creator client.ClientCreator) (res client.ClientCreator) {
	if creator == nil {
		res = func(address, checkId, output string) client.Client {
			return client.NewClient(address, checkId, output)
		}
	} else {
		res = creator
	}
	return
}

type semaphoreStategyOpts struct {
	Strategy func(*Service, Tactic) selectors.Semaphore
	Tactic   Tactic
}

func (o semaphoreStategyOpts) Apply(s *Service) selectors.Semaphore { return o.Strategy(s, o.Tactic) }

func ParseSemaphoreStrategy(opts StrategyOptions) (semaphoreStategyOpts, error) {
	var (
		strategy semaphoreStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseSemaphoreStrategy(opts.Strategy); err != nil {
		return semaphoreStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return semaphoreStategyOpts{}, err
	}

	return semaphoreStategyOpts{strategy, tactic}, nil
}

func parseSemaphoreStrategy(strategy string) (semaphoreStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopSemaphore, nil
	case "semaphore":
		return Semaphore, nil
	}
	return NoopSemaphore, typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid semaphore consul strategy %q", strategy)
}

type heartbeatStategyOpts struct {
	Strategy func(*Service, Tactic) selectors.Heartbeat
	Tactic   Tactic
}

func (o heartbeatStategyOpts) Apply(s *Service) selectors.Heartbeat { return o.Strategy(s, o.Tactic) }

func ParseHeartbeatStrategy(opts StrategyOptions) (heartbeatStategyOpts, error) {
	var (
		strategy heartbeatStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseHeartbeatStrategy(opts.Strategy); err != nil {
		return heartbeatStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return heartbeatStategyOpts{}, err
	}

	return heartbeatStategyOpts{strategy, tactic}, nil
}

func parseHeartbeatStrategy(strategy string) (heartbeatStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopHeartbeat, nil
	case "heartbeat":
		return Heartbeat, nil
	}
	return NoopHeartbeat, typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid heartbeat consul strategy %q", strategy)
}

type keyStoreStategyOpts struct {
	Strategy func(*Service, Tactic) selectors.KeyStore
	Tactic   Tactic
}

func (o keyStoreStategyOpts) Apply(s *Service) selectors.KeyStore { return o.Strategy(s, o.Tactic) }

func ParseKeyStoreStrategy(opts StrategyOptions) (keyStoreStategyOpts, error) {
	var (
		strategy keyStoreStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseKeyStoreStrategy(opts.Strategy); err != nil {
		return keyStoreStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return keyStoreStategyOpts{}, err
	}

	return keyStoreStategyOpts{strategy, tactic}, nil
}

func parseKeyStoreStrategy(strategy string) (keyStoreStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopKeyStore, nil
	case "keystore":
		return KeyStore, nil
	}
	return NoopKeyStore, typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid keyStore consul strategy %q", strategy)
}

func readTactic(tactic string,
	requestsPerDuration int,
	requestsDuration string,
) (Tactic, error) {
	dur, err := time.ParseDuration(requestsDuration)
	if err != nil {
		return noopTactic, err
	}

	switch common.Normalise(tactic) {
	case "noop":
		return noopTactic, nil
	case "nonblocking":
		return nonBlocking, nil
	case "ratelimited":
		return rateLimited(requestsPerDuration, dur), nil
	}
	return noopTactic, typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid consul tactic %q", tactic)
}
