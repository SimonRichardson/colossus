package store

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/permitters"
	r "github.com/SimonRichardson/colossus/cluster/store"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
)

func noopTactic([]r.Cluster, func(int, r.Cluster)) error {
	return nil
}

func nonBlocking(clusters []r.Cluster, fn func(int, r.Cluster)) error {
	for k, c := range clusters {
		go func(k int, c r.Cluster) {
			fn(k, c)
		}(k, c)
	}
	return nil
}

func rateLimited(requestsPerDuration int,
	requestsDuration time.Duration,
) func([]r.Cluster, func(int, r.Cluster)) error {
	permits := permitters.New(int64(requestsPerDuration), requestsDuration)
	return func(clusters []r.Cluster, fn func(int, r.Cluster)) error {
		if n := len(clusters); !permits.Allowed(int64(n)) {
			return typex.Errorf(errors.Source, errors.RateLimited, "element rate exceeded; request discarded")
		}
		return nonBlocking(clusters, fn)
	}
}
