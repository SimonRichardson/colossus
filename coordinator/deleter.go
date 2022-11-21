package coordinator

import (
	"github.com/SimonRichardson/colossus/farm/counter"
	"github.com/SimonRichardson/colossus/farm/store"
	s "github.com/SimonRichardson/colossus/selectors"
)

const (
	defaultRetryAmount = 3
)

type deleter struct {
	s.LifeCycleManager

	co      *Coordinator
	counter *counter.Farm
	store   *store.Farm
}

func newDeleter(co *Coordinator, counter *counter.Farm, store *store.Farm) *deleter {
	return &deleter{
		LifeCycleManager: newLifeCycleService(),

		co:      co,
		counter: counter,
		store:   store,
	}
}

func (i *deleter) Delete(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (int, error) {
	// Bucketize the members so that we can effiecently call all the storage
	// collections.
	var (
		instr   = i.co.instrumentation
		buckets = s.KeyFieldScoreTxnValues(members).Bucketize()
	)

	// We actually only care about this value, the rest is superfluous
	var (
		result         = 0
		partialFailure = false
	)
	for _, v := range buckets {
		res, err := i.store.Delete(v, maxSize)
		if err != nil {
			go instr.DeletePartialFailure()
			partialFailure = true
			continue
		}

		// The counter needs to retry as much as possible to
		updated := false
		for j := 0; j < defaultRetryAmount; j++ {
			if amount, err := i.counter.Delete(v, maxSize); err == nil && res == amount {
				updated = true
				break
			} else if res != amount {
				go instr.DeletePartialFailure()
			}

			// There was an error, so report it.
			partialFailure = true
		}

		// Don't report non-updated requests.
		if updated {
			result += res
		}
	}

	if partialFailure {
		return result, ErrPartialDeletionFailure
	}

	return result, nil
}

func (i *deleter) Rollback(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) error {
	var (
		instr   = i.co.instrumentation
		buckets = s.KeyFieldScoreTxnValues(members).Bucketize()

		partialFailure = false
	)

	for _, v := range buckets {
		if _, err := i.store.Delete(v, maxSize); err != nil {
			partialFailure = true
			go instr.RollbackPartialFailure()
		}

		if _, err := i.counter.Delete(v, maxSize); err != nil {
			partialFailure = true
			go instr.RollbackPartialFailure()
		}

		values := s.KeyFieldScoreTxnValues(v).KeyFieldScoreSizeExpiry(maxSize)
		if err := i.co.notifier.Unpublish(defaultInsertChannel, values); err != nil {
			partialFailure = true
			go instr.RollbackPartialFailure()
		}

		if err := i.co.persistence.Rollback(v, maxSize); err != nil {
			partialFailure = true
		}
	}

	if partialFailure {
		return ErrPartialRollbackFailure
	}
	return nil
}
