package cache

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
func DefaultConfig(address string, expiry time.Duration) ([]Cluster, error) {
	return ParseString(address, "30s", "10s", "10s", "hash", 1000, expiry, nil)
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
	expiry time.Duration,
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
			expiry,
		))
	}

	if len(clusters) < 1 {
		return empty, typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"No clusters specified %q", addresses)
	}

	return clusters, nil
}

type encodeStategyOpts struct {
	Strategy func(*Service, Tactic) s.Encoder
	Tactic   Tactic
}

func (o encodeStategyOpts) Apply(f *Service) s.Encoder { return o.Strategy(f, o.Tactic) }

func ParseEncodeStrategy(opts StrategyOptions) (encodeStategyOpts, error) {
	var (
		strategy encodeStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseEncodeStrategy(opts.Strategy, opts.Quorum); err != nil {
		return encodeStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return encodeStategyOpts{}, err
	}

	return encodeStategyOpts{strategy, tactic}, nil
}

func parseEncodeStrategy(strategy string, quorum float64) (encodeStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopEncoder, nil
	case "encoder":
		return Encoder, nil
	}
	return NoopEncoder, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid encode redis strategy %q", strategy)
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
