package lorenz

import (
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/SimonRichardson/colossus/blist/common"
	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/blist/services/lorenz/client"
	"github.com/SimonRichardson/colossus/typex"
)

// DefaultConfig creates a default configuration for the service
func DefaultConfig(addresses, version string) ([]Cluster, error) {
	return ParseString(addresses, version, "30s", nil)
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
func ParseString(addresses, version, timeout string, creator client.ClientCreator) ([]Cluster, error) {
	var (
		clusters      = []Cluster{}
		fnCreator     = clientCreator(creator)
		timeouts, err = ParseTimeout(timeout)
	)

	if err != nil {
		return empty, err
	}

	for _, address := range strings.Split(common.StripWhitespace(addresses), ";") {
		u, err := url.Parse(address)
		if err != nil {
			return empty, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
				"Invalid url %q for lorenz (%s)", address, err)
		}

		if strings.Contains(u.Host, ":") {
			host, port, err := net.SplitHostPort(u.Host)
			if err != nil {
				return empty, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
					"Invalid host %q and port %q for lorenz (%s)", host, port, err)
			}

			if _, err := strconv.ParseUint(port, 10, 16); err != nil {
				return empty, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
					"Invalid port %q in host %q for lorenz (%s)", port, host, err)
			}
		}

		clusters = append(clusters, newCluster(
			fnCreator(address, version, timeouts),
		))
	}

	if len(clusters) < 1 {
		return empty, typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"No clusters specified %q for lorenz", addresses)
	}

	return clusters, nil
}

func clientCreator(creator client.ClientCreator) (res client.ClientCreator) {
	if creator == nil {
		res = func(address, version string, timeout *client.ConnectionTimeout) client.Client {
			return client.New(address, version, timeout)
		}
	} else {
		res = creator
	}
	return
}

// Parse defines away to create a client from a set of strings.
func ParseTimeout(global string) (*client.ConnectionTimeout, error) {
	timeout := client.NewConnectionTimeout()
	if dur, err := time.ParseDuration(global); err != nil {
		return nil, typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"Invalid timeout parssed %q", global)
	} else {
		timeout.Global = dur
	}
	return timeout, nil
}

type chargeStategyOpts struct {
	Strategy func(*Service, Tactic) selectors.Charger
	Tactic   Tactic
}

func (o chargeStategyOpts) Apply(s *Service) selectors.Charger { return o.Strategy(s, o.Tactic) }

type eventSelectStategyOpts struct {
	Strategy func(*Service, Tactic) selectors.EventSelector
	Tactic   Tactic
}

func (o eventSelectStategyOpts) Apply(s *Service) selectors.EventSelector {
	return o.Strategy(s, o.Tactic)
}

type codeSetSelectStategyOpts struct {
	Strategy func(*Service, Tactic) selectors.CodeSetSelector
	Tactic   Tactic
}

func (o codeSetSelectStategyOpts) Apply(s *Service) selectors.CodeSetSelector {
	return o.Strategy(s, o.Tactic)
}

type inspectStategyOpts struct {
	Strategy func(*Service, Tactic) selectors.Inspector
	Tactic   Tactic
}

func (o inspectStategyOpts) Apply(s *Service) selectors.Inspector { return o.Strategy(s, o.Tactic) }

func ParseChargeStrategy(opts StrategyOptions) (chargeStategyOpts, error) {
	var (
		strategy chargeStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseChargeStrategy(opts.Strategy); err != nil {
		return chargeStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return chargeStategyOpts{}, err
	}

	return chargeStategyOpts{strategy, tactic}, nil
}

func parseChargeStrategy(strategy string) (chargeStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopCharger, nil
	case "charger":
		return Charger, nil
	}
	return NoopCharger, typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid charger lorenz strategy %q", strategy)
}

func ParseEventSelectStrategy(opts StrategyOptions) (eventSelectStategyOpts, error) {
	var (
		strategy eventSelectorStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseEventSelectStrategy(opts.Strategy); err != nil {
		return eventSelectStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return eventSelectStategyOpts{}, err
	}

	return eventSelectStategyOpts{strategy, tactic}, nil
}

func parseEventSelectStrategy(strategy string) (eventSelectorStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopEventSelector, nil
	case "eventselector":
		return EventSelector, nil
	case "eventcacheselector":
		return EventCacheSelector, nil
	}
	return NoopEventSelector, typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid event select lorenz strategy %q", strategy)
}

func ParseCodeSetSelectStrategy(opts StrategyOptions) (codeSetSelectStategyOpts, error) {
	var (
		strategy codeSetSelectorStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseCodeSetSelectStrategy(opts.Strategy); err != nil {
		return codeSetSelectStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return codeSetSelectStategyOpts{}, err
	}

	return codeSetSelectStategyOpts{strategy, tactic}, nil
}

func parseCodeSetSelectStrategy(strategy string) (codeSetSelectorStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopCodeSetSelector, nil
	case "codesetselector":
		return CodeSetSelector, nil
	}
	return NoopCodeSetSelector, typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid codeset select lorenz strategy %q", strategy)
}

func ParseInspectStrategy(opts StrategyOptions) (inspectStategyOpts, error) {
	var (
		strategy inspectStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseInspectStrategy(opts.Strategy); err != nil {
		return inspectStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return inspectStategyOpts{}, err
	}

	return inspectStategyOpts{strategy, tactic}, nil
}

func parseInspectStrategy(strategy string) (inspectStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopInspector, nil
	case "inspector":
		return Inspector, nil
	}
	return NoopInspector, typex.Errorf(errors.Source, errors.UnexpectedArgument,
		"Invalid inspect lorenz strategy %q", strategy)
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
		"Invalid lorenz tactic %q", tactic)
}
