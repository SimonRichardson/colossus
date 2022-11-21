package redis

import (
	"time"

	"github.com/SimonRichardson/colossus/fusion"
	"github.com/SimonRichardson/colossus/fusion/common"
	"github.com/SimonRichardson/colossus/fusion/errors"
	"github.com/SimonRichardson/colossus/fusion/selectors"
	"github.com/SimonRichardson/colossus/fusion/strategies"
	"github.com/SimonRichardson/colossus/typex"
)

type RedisOptions struct {
	KeyStorePrefix selectors.Prefix
	KeyStoreTicker chan struct{}
	KeyStore       selectors.KeyStore
}

// Parse defines away to create a redis pool from a set of strings.
func Parse(connect, read, write, strategy string,
	opts *RedisOptions,
) (*ConnectionTimeout, fusion.SelectionStrategy, error) {
	items := map[string]string{
		"connect": connect,
		"read":    read,
		"write":   write,
	}

	if timeout, err := readTimeout(items); err != nil {
		return nil, nil, err
	} else if selection, err := readStrategy(strategy, opts); err != nil {
		return nil, nil, err
	} else {
		return timeout, selection, nil
	}
}

func readTimeout(items map[string]string) (*ConnectionTimeout, error) {
	timeout := newConnectionTimeout()
	for k, v := range items {
		dur, err := time.ParseDuration(v)
		if err != nil {
			return nil, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
				"Invalid timeout %q passed to %s", v, k)
		}

		switch k {
		case "connect":
			timeout.connect = dur
		case "read":
			timeout.read = dur
		case "write":
			timeout.write = dur
		}
	}
	return timeout, nil
}

func readStrategy(strategy string, opts *RedisOptions) (fusion.SelectionStrategy, error) {
	switch common.Normalise(strategy) {
	case "hash":
		return strategies.NewHash(), nil
	case "roundrobin":
		return strategies.NewRoundRobin(), nil
	case "random":
		return strategies.NewRandom(), nil
	case "keystore":
		if opts == nil {
			return nil, typex.Errorf(errors.Source, errors.UnexpectedArgument,
				"Missing key store options for strategy %q", strategy)
		}
		return strategies.NewKeyStore(
			opts.KeyStorePrefix,
			opts.KeyStoreTicker,
			opts.KeyStore,
		), nil
	}
	return nil, typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
		"Invalid pool selection strategy %q", strategy)
}
