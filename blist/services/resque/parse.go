package resque

import (
	"strings"
	"time"

	"github.com/SimonRichardson/colossus/blist/common"
	"github.com/SimonRichardson/colossus/blist/errors"
	s "github.com/SimonRichardson/colossus/blist/selectors"
	r "github.com/SimonRichardson/colossus/fusion/redis"
	"github.com/SimonRichardson/colossus/typex"
)

// DefaultConfig creates a default configuration for the service
func DefaultConfig(address string) ([]Cluster, error) {
	return ParseString(address, "1m", "10s", "20s", "Hash", 1000, nil)
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
// then use with in the farm.
//   - addresses is a semi-colon separated string of redis addresses
//   - connectTimeout, readTimeout and writeTimeout is a set of durations in
//     string format
//   - poolRoutingStrategy defines a strategy for how the pool routing works
func ParseString(addresses string,
	connectTimeout, readTimeout, writeTimeout string,
	poolRoutingStrategy string,
	maxSize int,
	creator r.RedisCreator,
) ([]Cluster, error) {
	var (
		clusters                = []Cluster{}
		timeouts, strategy, err = r.Parse(connectTimeout,
			readTimeout,
			writeTimeout,
			poolRoutingStrategy,
			nil,
		)
	)

	if err != nil {
		return empty, err
	}

	for i, address := range strings.Split(common.StripWhitespace(addresses), ";") {
		hosts := []string{}
		for _, host := range strings.Split(address, ",") {
			if host == "" {
				continue
			}
			if err := r.ValidRedisHost(host); err != nil {
				return empty, err
			}
			hosts = append(hosts, host)
		}

		if len(hosts) < 1 {
			return empty, typex.Errorf(errors.Source, errors.UnexpectedArgument,
				"Empty cluster %d (%q)", i+1, address)
		}

		clusters = append(clusters, newCluster(
			r.New(hosts, strategy, timeouts, maxSize, creator),
		))
	}

	if len(clusters) < 1 {
		return empty, typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"No clusters specified %q", addresses)
	}

	return clusters, nil
}

type enqueueStategyOpts struct {
	Strategy func(*Service, Tactic) s.Enqueuer
	Tactic   Tactic
}

func (o enqueueStategyOpts) Apply(f *Service) s.Enqueuer { return o.Strategy(f, o.Tactic) }

func ParseEnqueueStrategy(opts StrategyOptions) (enqueueStategyOpts, error) {
	var (
		strategy enqueueStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseEnqueueStrategy(opts.Strategy, opts.Quorum); err != nil {
		return enqueueStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return enqueueStategyOpts{}, err
	}

	return enqueueStategyOpts{strategy, tactic}, nil
}

func parseEnqueueStrategy(strategy string, quorum float64) (enqueueStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopEnqueuer, nil
	case "enqueuer":
		return Enqueuer, nil
	}
	return NoopEnqueuer, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid enqueue redis strategy %q", strategy)
}

type registerStategyOpts struct {
	Strategy func(*Service, Tactic) s.Register
	Tactic   Tactic
}

func (o registerStategyOpts) Apply(f *Service) s.Register { return o.Strategy(f, o.Tactic) }

func ParseRegisterStrategy(opts StrategyOptions) (registerStategyOpts, error) {
	var (
		strategy registerStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseRegisterStrategy(opts.Strategy, opts.Quorum); err != nil {
		return registerStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return registerStategyOpts{}, err
	}

	return registerStategyOpts{strategy, tactic}, nil
}

func parseRegisterStrategy(strategy string, quorum float64) (registerStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopRegister, nil
	case "register":
		return Register, nil
	}
	return NoopRegister, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid register redis strategy %q", strategy)
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
	return noopTactic, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid redis tactic %q", tactic)
}
