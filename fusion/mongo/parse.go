package mongo

import (
	"time"

	"github.com/SimonRichardson/colossus/fusion"
	"github.com/SimonRichardson/colossus/fusion/common"
	"github.com/SimonRichardson/colossus/fusion/errors"
	"github.com/SimonRichardson/colossus/fusion/strategies"
	"github.com/SimonRichardson/colossus/typex"
)

// Parse defines away to create a mongo pool from a set of strings.
func Parse(global, strategy string) (*ConnectionTimeout, fusion.SelectionStrategy, error) {
	if timeout, err := readTimeout(global); err != nil {
		return nil, nil, err
	} else if selection, err := readStrategy(strategy); err != nil {
		return nil, nil, err
	} else {
		return timeout, selection, nil
	}
}

func readTimeout(global string) (*ConnectionTimeout, error) {
	timeout := newConnectionTimeout()
	if dur, err := time.ParseDuration(global); err != nil {
		return nil, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
			"Invalid timeout passed %s", global)
	} else {
		timeout.global = dur
	}
	return timeout, nil
}

func readStrategy(strategy string) (fusion.SelectionStrategy, error) {
	switch common.Normalise(strategy) {
	case "hash":
		return strategies.NewHash(), nil
	case "roundrobin":
		return strategies.NewRoundRobin(), nil
	case "random":
		return strategies.NewRandom(), nil
	}
	return nil, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid pool selection strategy %q", strategy)
}
