package notifier

import (
	"strings"
	"time"

	c "github.com/SimonRichardson/colossus/cluster/notifier"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/errors"
	r "github.com/SimonRichardson/colossus/fusion/redis"
	s "github.com/SimonRichardson/colossus/selectors"
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
			if len(host) < 1 {
				continue
			}
			if err := r.ValidRedisHost(host); err != nil {
				return empty, err
			}
			hosts = append(hosts, host)
		}

		if len(hosts) < 1 {
			return empty, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
				"Empty cluster %d (%q)", i+1, address)
		}

		clusters = append(clusters, c.New(
			r.New(hosts, strategy, timeouts, maxSize, creator),
		))
	}

	if len(clusters) < 1 {
		return empty, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
			"No clusters specified %q", addresses)
	}

	return clusters, nil
}

type notifyStategyOpts struct {
	Strategy func(*Farm, Tactic) s.Notifier
	Tactic   Tactic
}

func (o notifyStategyOpts) Apply(f *Farm) s.Notifier { return o.Strategy(f, o.Tactic) }

func ParseNotifyStrategy(opts env.StrategyOptions) (notifyStategyOpts, error) {
	var (
		strategy notifierStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseNotifyStrategy(opts.Strategy, opts.Quorum); err != nil {
		return notifyStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return notifyStategyOpts{}, err
	}

	return notifyStategyOpts{strategy, tactic}, nil
}

func parseNotifyStrategy(strategy string, quorum float64) (notifierStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopNotifier, nil
	case "individual":
		return Individual, nil
	case "bulk":
		return Bulk, nil
	}
	return NoopNotifier, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid insert notifier strategy %q", strategy)
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
		"Invalid notifier tactic %q", tactic)
}
