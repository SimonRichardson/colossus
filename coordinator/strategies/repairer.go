package strategies

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/permitters"
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/farm/store"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// RepairStrategy defines if a repair can actually happen or if it'll be refused
type RepairStrategy func(*store.Farm, []s.KeyFieldTxnValue, s.KeySizeExpiry) error

func repairNonBlocking(farm *store.Farm, members []s.KeyFieldTxnValue, maxSize s.KeySizeExpiry) error {
	buckets := map[bs.Key][]s.KeyFieldTxnValue{}
	for _, v := range members {
		buckets[v.Key] = append(buckets[v.Key], v)
	}

	responses := make(chan error, len(buckets))

	wg := sync.WaitGroup{}
	wg.Add(len(buckets))
	go func() { wg.Wait(); close(responses) }()

	for _, v := range buckets {
		go func(elements []s.KeyFieldTxnValue) {
			if err := farm.Repair(elements, maxSize); err != nil {
				responses <- err
			}
		}(v)
	}

	errs := []error{}
	for e := range responses {
		errs = append(errs, e)
	}

	return typex.Errorf(errors.Source, errors.NoCaseFound,
		"Error Repairing (%s)", common.SumErrors(errs).Error())
}

func repairNoopTactic(*store.Farm, []s.KeyFieldTxnValue, s.KeySizeExpiry) error {
	return nil
}

func repairRateLimited(maxElements int64, maxDuration time.Duration, strategy RepairStrategy) RepairStrategy {
	permits := permitters.New(maxElements, maxDuration)
	return func(farm *store.Farm, members []s.KeyFieldTxnValue, maxSize s.KeySizeExpiry) error {
		if num := len(members); !permits.Allowed(int64(num)) {
			return nil
		}

		return strategy(farm, members, maxSize)
	}
}
