package strategies

import (
	"time"

	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
)

func NewManagerStrategy(e *env.Env) (ManagerStrategyCreator, error) {
	options := e.GetRepairOptions(env.Manager)

	switch common.Normalise(options.Strategy) {
	case "collect":
		dur, err := time.ParseDuration(options.RequestsDuration)
		if err != nil {
			return managerNoopStrategy, err
		}
		return managerCollectStrategy(dur), nil
	case "noop":
		return managerNoopStrategy, nil
	}
	return managerNoopStrategy, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"No strategy found")
}

func NewRepairStrategy(e *env.Env) (RepairStrategy, error) {
	options := e.GetRepairOptions(env.Coordinator)

	switch common.Normalise(options.Tactic) {
	case "ratelimited":
		var (
			maxElements = options.RequestsPerDuration
			dur, err    = time.ParseDuration(options.RequestsDuration)
		)
		if err != nil {
			return repairNoopTactic, err
		}
		stategy, err := NewRepairStrategy(e)
		if err != nil {
			return repairNoopTactic, err
		}
		return repairRateLimited(int64(maxElements), dur, stategy), nil
	case "nonblocking":
		return repairNonBlocking, nil
	case "noop":
		return repairNoopTactic, nil
	}

	return repairNoopTactic, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"No tactic found")
}

func NewInsertStrategy(e *env.Env) (InsertStrategy, error) {
	options := e.GetInsertOptions(env.Coordinator)

	switch common.Normalise(options.Tactic) {
	case "counter":
		return insertCounter, nil
	case "noop":
		return insertNoop, nil
	}

	return insertNoop, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"No tactic found")
}
