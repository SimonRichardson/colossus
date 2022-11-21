package store

import (
	"strings"
	"time"

	c "github.com/SimonRichardson/colossus/cluster/store"
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
	opts *r.RedisOptions,
	creator r.RedisCreator,
) ([]c.Cluster, error) {
	var (
		clusters                = []c.Cluster{}
		timeouts, strategy, err = r.Parse(connectTimeout,
			readTimeout,
			writeTimeout,
			poolRoutingStrategy,
			opts,
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

type selectStategyOpts struct {
	Strategy func(*Farm, Tactic) s.Selector
	Tactic   Tactic
}

func (o selectStategyOpts) Apply(f *Farm) s.Selector { return o.Strategy(f, o.Tactic) }

type insertStategyOpts struct {
	Strategy func(*Farm, Tactic) s.Inserter
	Tactic   Tactic
}

func (o insertStategyOpts) Apply(f *Farm) s.Inserter { return o.Strategy(f, o.Tactic) }

type modifyStategyOpts struct {
	Strategy func(*Farm, Tactic) s.Modifier
	Tactic   Tactic
}

func (o modifyStategyOpts) Apply(f *Farm) s.Modifier { return o.Strategy(f, o.Tactic) }

type deleteStategyOpts struct {
	Strategy func(*Farm, Tactic) s.Deleter
	Tactic   Tactic
}

func (o deleteStategyOpts) Apply(f *Farm) s.Deleter { return o.Strategy(f, o.Tactic) }

type scanStategyOpts struct {
	Strategy func(*Farm, Tactic) s.Scanner
	Tactic   Tactic
}

func (o scanStategyOpts) Apply(f *Farm) s.Scanner { return o.Strategy(f, o.Tactic) }

type repairStategyOpts struct {
	Strategy func(*Farm, Tactic) s.Repairer
	Tactic   Tactic
}

func (o repairStategyOpts) Apply(f *Farm) s.Repairer { return o.Strategy(f, o.Tactic) }

func ParseSelectStrategy(opts env.StrategyOptions) (selectStategyOpts, error) {
	var (
		strategy selectStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseSelectStrategy(opts.Strategy, opts.Quorum); err != nil {
		return selectStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return selectStategyOpts{}, err
	}

	return selectStategyOpts{strategy, tactic}, nil
}

func ParseInsertStrategy(opts env.StrategyOptions) (insertStategyOpts, error) {
	var (
		strategy insertStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseInsertStrategy(opts.Strategy, opts.Quorum); err != nil {
		return insertStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return insertStategyOpts{}, err
	}

	return insertStategyOpts{strategy, tactic}, nil
}

func ParseDeleteStrategy(opts env.StrategyOptions) (deleteStategyOpts, error) {
	var (
		strategy deleteStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseDeleteStrategy(opts.Strategy, opts.Quorum); err != nil {
		return deleteStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return deleteStategyOpts{}, err
	}

	return deleteStategyOpts{strategy, tactic}, nil
}

func ParseScanStrategy(opts env.StrategyOptions) (scanStategyOpts, error) {
	var (
		strategy scanStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseScanStrategy(opts.Strategy, opts.Quorum); err != nil {
		return scanStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return scanStategyOpts{}, err
	}

	return scanStategyOpts{strategy, tactic}, nil
}

func ParseRepairStrategy(opts env.StrategyOptions) (repairStategyOpts, error) {
	var (
		strategy repairStrategy
		tactic   Tactic
		err      error
	)

	if strategy, err = parseRepairStrategy(opts.Strategy, opts.Quorum); err != nil {
		return repairStategyOpts{}, err
	}
	if tactic, err = readTactic(opts.Tactic,
		opts.RequestsPerDuration,
		opts.RequestsDuration,
	); err != nil {
		return repairStategyOpts{}, err
	}

	return repairStategyOpts{strategy, tactic}, nil
}

func parseSelectStrategy(strategy string, quorum float64) (selectStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopSelector, nil
	case "selectonereadone":
		return SelectOneReadOne, nil
	case "selectallreadall":
		return SelectAllReadAll, nil
	case "selectquorumreadall":
		return SelectQuorumReadAll(quorum), nil
	}
	return NoopSelector, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid store strategy %q", strategy)
}

func parseInsertStrategy(strategy string, quorum float64) (insertStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopInserter, nil
	case "insertallreadall":
		return InsertAllReadAll, nil
	case "insertallreadquorum":
		return InsertAllReadQuorum(quorum), nil
	}
	return NoopInserter, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid insert store strategy %q", strategy)
}

func parseDeleteStrategy(strategy string, quorum float64) (deleteStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopDeleter, nil
	case "deleteallreadall":
		return DeleteAllReadAll, nil
	case "deleteallreadquorum":
		return DeleteAllReadQuorum(quorum), nil
	}
	return NoopDeleter, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid delete store strategy %q", strategy)
}

func parseScanStrategy(strategy string, quorum float64) (scanStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopScanner, nil
	case "scanallreadall":
		return ScanAllReadAll, nil
	}
	return NoopScanner, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid store strategy %q", strategy)
}

func parseRepairStrategy(strategy string, quorum float64) (repairStrategy, error) {
	switch common.Normalise(strategy) {
	case "noop":
		return NoopRepairer, nil
	case "repairall":
		return RepairAll, nil
	}
	return NoopRepairer, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid store strategy %q", strategy)
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
		"Invalid store tactic %q", tactic)
}
