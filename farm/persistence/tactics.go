package persistence

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/permitters"
	p "github.com/SimonRichardson/colossus/cluster/persistence"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
)

func noopTactic([]p.Cluster, func(p.Cluster)) error {
	return nil
}

func nonBlocking(clusters []p.Cluster, fn func(p.Cluster)) error {
	for _, c := range clusters {
		go func(c p.Cluster) {
			fn(c)
		}(c)
	}
	return nil
}

func rateLimited(requestsPerDuration int,
	requestsDuration time.Duration,
) func([]p.Cluster, func(p.Cluster)) error {
	permits := permitters.New(int64(requestsPerDuration), requestsDuration)
	return func(clusters []p.Cluster, fn func(p.Cluster)) error {
		if n := len(clusters); !permits.Allowed(int64(n)) {
			return typex.Errorf(errors.Source, errors.RateLimited, "Element rate exceeded; repair request discarded")
		}
		return nonBlocking(clusters, fn)
	}
}
