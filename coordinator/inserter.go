package coordinator

import (
	"github.com/SimonRichardson/colossus/coordinator/strategies"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/farm/counter"
	"github.com/SimonRichardson/colossus/farm/notifier"
	"github.com/SimonRichardson/colossus/farm/store"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/typex"
)

type inserter struct {
	s.LifeCycleManager

	co       *Coordinator
	counter  *counter.Farm
	store    *store.Farm
	notifier *notifier.Farm
	strategy strategies.InsertStrategy
}

func newInserter(co *Coordinator,
	counter *counter.Farm,
	store *store.Farm,
	notifier *notifier.Farm,
	strategy strategies.InsertStrategy,
) *inserter {
	return &inserter{
		LifeCycleManager: newLifeCycleService(),

		co:       co,
		counter:  counter,
		store:    store,
		notifier: notifier,
		strategy: strategy,
	}
}

func (i *inserter) Insert(members []s.KeyFieldScoreTxnValue, sizeExpiry s.KeySizeExpiry) (int, error) {
	// Bucketize the members so that we can effiecently call all the storage
	// collections.
	var (
		instr      = i.co.instrumentation
		buckets    = s.KeyFieldScoreTxnValues(members).Bucketize()
		sized, err = i.strategy(i.counter, buckets, sizeExpiry)
	)

	if err != nil {
		return 0, err
	}

	// We actually only care about this value, the rest is superfluous
	var (
		result         = 0
		partialFailure = false
	)
	for k, v := range sized {
		res, err := i.store.Insert(v, sizeExpiry)
		if err != nil {

			teleprinter.L.Error().Printf("Store Insert Partial Failure (%s, %d:%d)",
				k.String(), len(v), res)

			go instr.InsertPartialFailure()
			return result, typex.Errorf(errors.Source, errors.Partial,
				"Partial insertion failure (%s, %s)", k.String(), err.Error())
		} else if num := len(v); res != num {
			// If we reach here, then we *may* of inserted some items already
			// and we're in a situation that we should rollback previous
			// inserted ones.
			teleprinter.L.Error().Printf("Store Insert Partial Failure (%s, %d:%d)",
				k.String(), len(v), res)

			go instr.InsertPartialFailure()
			return result, typex.Errorf(errors.Source, errors.Partial,
				"Partial insertion failure (%s, %d:%d)", k.String(), res, num)
		}

		updated := false
		for j := 0; j < defaultRetryAmount; j++ {
			amount, err := i.counter.Insert(v, sizeExpiry)
			if err == nil && res == amount {
				updated = true
				break
			} else if res != amount {
				go instr.InsertPartialFailure()
			}

			// There was an error, so report it
			partialFailure = true
			teleprinter.L.Error().Printf("Counter Insert Partial Failure (%s, %d:%d)\n",
				k.String(), res, amount)
		}

		// Don't report non-updated requests.
		if updated {
			result += res
		}
	}

	if partialFailure {
		teleprinter.L.Error().Printf("Counter Failed Insertion with Partial Error\n")

		return result, ErrPartialInsertionFailure
	}

	go i.notifier.Publish(defaultInsertChannel, s.KeyFieldScoreTxnValues(members).KeyFieldScoreSizeExpiry(sizeExpiry))

	return result, nil
}
