package cache

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/permitters"
	"github.com/SimonRichardson/colossus/typex"
)

func noopTactic([]Cluster, func(Cluster)) error {
	return nil
}

func nonBlocking(clusters []Cluster, fn func(Cluster)) error {
	for _, c := range clusters {
		go func(c Cluster) {
			fn(c)
		}(c)
	}
	return nil
}

func rateLimited(requestPerDuration int,
	requestsDuration time.Duration,
) func([]Cluster, func(Cluster)) error {
	permits := permitters.New(int64(requestPerDuration), requestsDuration)
	return func(clusters []Cluster, fn func(Cluster)) error {
		if n := len(clusters); !permits.Allowed(int64(n)) {
			return typex.Errorf(errors.Source, errors.RateLimited,
				"RateLimited: element rate exceeded; request discarded")
		}
		return nonBlocking(clusters, fn)
	}
}
