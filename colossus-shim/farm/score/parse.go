package score

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	c "github.com/SimonRichardson/colossus/colossus-shim/cluster/score"
	t "github.com/SimonRichardson/colossus/colossus-shim/selectors"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/errors"
	r "github.com/SimonRichardson/colossus/fusion/redis"
	"github.com/SimonRichardson/colossus/typex"
)

var (
	empty = []c.Cluster{}
)

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
) ([]c.Cluster, error) {
	var (
		clusters                = []c.Cluster{}
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
			if strings.Contains(host, ":") {
				url, err := url.Parse(host)
				if err != nil {
					return empty, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
						"Invalid host %q (%s)", host, err)
				}

				tokens := strings.Split(url.Host, ":")
				if len(tokens) < 2 {
					return empty, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
						"Invalid host %q", host)
				}
				if _, err := strconv.ParseUint(tokens[1], 10, 16); err != nil {
					return empty, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
						"Invalid port %q in host %q (%s)", tokens[1], host, err)
				}

				if url.User != nil {
					if password, ok := url.User.Password(); ok && len(password) < 1 {
						return empty, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
							"Invalid password %q in host %q", password, host)
					}
				}
			}

			hosts = append(hosts, host)
		}

		if len(hosts) < 1 {
			return empty, typex.Errorf(errors.Source, errors.UnexpectedArgument,
				"Empty cluster %d (%q)", i+1, address)
		}

		clusters = append(clusters, c.New(
			r.New(hosts, strategy, timeouts, maxSize, creator),
		))
	}

	if len(clusters) < 1 {
		return empty, typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"No clusters specified %q", addresses)
	}

	return clusters, nil
}

type incrementStategyOpts struct {
	Strategy func(*Farm, Tactic) t.Incrementer
	Tactic   Tactic
}

func (o incrementStategyOpts) Apply(f *Farm) t.Incrementer { return o.Strategy(f, o.Tactic) }

func ParseIncrementStrategy(opts env.StrategyOptions) (incrementStategyOpts, error) {
	var (
		strategy incrementStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseIncrementStrategy(opts.Strategy, opts.Quorum); err != nil {
		return incrementStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return incrementStategyOpts{}, err
	}

	return incrementStategyOpts{strategy, tactic}, nil
}

func parseIncrementStrategy(strategy string, quorum float64) (incrementStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopIncrementer, nil
	case "incrementallreadall":
		return IncrementAllReadAll, nil
	case "incrementfromtime":
		return IncrementFromTime, nil
	}
	return NoopIncrementer, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid increment redis strategy %q", strategy)
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
