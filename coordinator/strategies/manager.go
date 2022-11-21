package strategies

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	cs "github.com/SimonRichardson/colossus/cluster/store"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/farm/store"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

const (
	defaultManagerSize = 1000
)

// ErrFatal defines if the manager should be shut down or not
var (
	ErrFatal = typex.Errorf(errors.Source, errors.Fatal, "Fatal error")

	defaultStep          = time.Minute.Nanoseconds()
	defaultIntervalSweep = time.Minute.Nanoseconds()
	defaultFullSweep     = (time.Minute * 10).Nanoseconds()
)

type Manager interface {
	s.Deleter
	s.Scanner
}

// ManagerStrategyCreator creates a ManagerStrategy
type ManagerStrategyCreator func(Manager, *store.Farm) ManagerStrategy

// ManagerStrategy defines how
type ManagerStrategy func(s.KeyFieldScoreSizeExpiry) error

func managerNoopStrategy(Manager, *store.Farm) ManagerStrategy {
	return func(kfs s.KeyFieldScoreSizeExpiry) error {
		return nil
	}
}

func managerCollectStrategy(duration time.Duration) ManagerStrategyCreator {
	return func(co Manager, sf *store.Farm) ManagerStrategy {
		var (
			// Run slightly head of speed so we don't miss time anything
			percentage    = int64(float64(defaultIntervalSweep) * 0.9)
			intervalTimer = time.NewTicker(time.Duration(percentage))
			fullTimer     = time.NewTicker(time.Duration(defaultFullSweep))

			slots = newTimeSlots()
		)

		go func() {
			for {
				select {
				case <-intervalTimer.C:
					// Check to see if the item has expired, if it has delete it!
					var (
						// Attempt to make the score be in the future
						items = slots.Peek()

						now    = time.Now()
						values = []s.KeyFieldScoreTxnValue{}
					)

					for _, v := range items {
						if item, ok := selectItem(sf, now, v.Key, v.Field); ok {
							values = append(values, item)
						}
					}

					if _, err := co.Delete(values, s.MakeKeySizeExpiry()); err != nil {
						log.Println("Partial failure", err)
					}

				case <-fullTimer.C:
					// Get all the keys then all the fields and then check to
					// see if the item has expired, if it has delete it!

					// TODO : How do we know another colossus isn't doing this
					// at the same time?
					keys, err := co.Keys()
					if err != nil {
						continue
					}

					shuffle(keys)

					var (
						now    = time.Now()
						values = []s.KeyFieldScoreTxnValue{}
					)

					for _, key := range keys {
						fields, err := co.Members(key)
						if err != nil {
							continue
						}

						for _, field := range fields {
							if item, ok := selectItem(sf, now, key, field); ok {
								values = append(values, item)
							}
						}
					}

					if _, err := co.Delete(values, s.MakeKeySizeExpiry()); err != nil {
						log.Println("Partial failure", err)
					}
				}
			}
		}()

		return func(kfs s.KeyFieldScoreSizeExpiry) error {
			return slots.Add(kfs)
		}
	}
}

func selectItem(sf *store.Farm, now time.Time, key, field bs.Key) (s.KeyFieldScoreTxnValue, bool) {
	item, err := sf.Select(key, field)
	if err != nil {
		if err != cs.ErrExpiredNode {
			return s.KeyFieldScoreTxnValue{}, false
		}
	}

	return s.KeyFieldScoreTxnValue{
		Key:   item.Key,
		Field: item.Field,
		Score: item.Score + 1,
		Txn:   item.Txn,
		Value: item.Value,
	}, true
}

type timeRange struct {
	start, end int64
}

func makeTimeRange(interval int64, offset int64) timeRange {
	return timeRange{interval, interval + offset}
}

type timeSlots struct {
	mutex *sync.Mutex
	slots map[timeRange][]s.KeyFieldScoreSizeExpiry
}

func newTimeSlots() *timeSlots {
	return &timeSlots{
		&sync.Mutex{},
		map[timeRange][]s.KeyFieldScoreSizeExpiry{},
	}
}

func (b *timeSlots) Add(k s.KeyFieldScoreSizeExpiry) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var (
		now  = time.Now()
		step = stepRound(now.Add(k.Expiry).UnixNano(), defaultStep)
		slot = makeTimeRange(step, defaultIntervalSweep)
	)
	b.slots[slot] = append(b.slots[slot], k)

	return nil
}

func (b *timeSlots) Peek() []s.KeyFieldScoreSizeExpiry {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var (
		step = stepRound(time.Now().UnixNano(), defaultStep)
		slot = makeTimeRange(step, defaultIntervalSweep)
	)

	if slots, ok := b.slots[slot]; ok {
		delete(b.slots, slot)
		return slots
	}

	return []s.KeyFieldScoreSizeExpiry{}
}

func stepRound(v int64, d int64) int64 {
	return int64(math.Ceil(float64(v/d))) * d
}

func shuffle(keys []bs.Key) {
	for i := range keys {
		j := rand.Intn(i + 1)
		keys[i], keys[j] = keys[j], keys[i]
	}
}
