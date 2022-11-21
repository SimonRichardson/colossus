package notifier

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/permitters"
	n "github.com/SimonRichardson/colossus/cluster/notifier"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
)

func noopTactic([]n.Cluster, func(n.Cluster)) error {
	return nil
}

func nonBlocking(clusters []n.Cluster, fn func(n.Cluster)) error {
	for _, c := range clusters {
		go func(c n.Cluster) {
			fn(c)
		}(c)
	}
	return nil
}

func rateLimited(requestsPerDuration int,
	requestsDuration time.Duration,
) func([]n.Cluster, func(n.Cluster)) error {
	permits := permitters.New(int64(requestsPerDuration), requestsDuration)
	return func(clusters []n.Cluster, fn func(n.Cluster)) error {
		if n := len(clusters); !permits.Allowed(int64(n)) {
			return typex.Errorf(errors.Source, errors.RateLimited,
				"element rate exceeded; repair request discarded")
		}
		return nonBlocking(clusters, fn)
	}
}
