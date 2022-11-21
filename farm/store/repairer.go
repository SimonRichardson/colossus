package store

import (
	"log"
	"strings"
	"sync"
	"time"

	t "github.com/SimonRichardson/colossus/cluster"
	"github.com/SimonRichardson/colossus/errors"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// RepairAll defines a strategy for attempting to repair all the nodes via a
// strong consensus.
func RepairAll(f *Farm, t Tactic) s.Repairer {
	return repair{f, t}
}

type repair struct {
	*Farm
	tactic Tactic
}

func (w repair) Repair(keyFieldTxnValue []s.KeyFieldTxnValue, maxSize s.KeySizeExpiry) error {
	var (
		clusters      = w.Farm.clusters
		numOfClusters = len(clusters)
		began         = time.Now()
	)

	go func() {
		instr := w.Farm.instrumentation
		instr.RepairCall()
		instr.RepairRequest(len(keyFieldTxnValue))
		instr.RepairSendTo(numOfClusters)
	}()
	defer func() {
		go func() {
			instr := w.Farm.instrumentation
			instr.RepairDuration(time.Since(began))
		}()
	}()

	presenceMap := map[s.KeyFieldTxnValue][]s.Presence{}
	for _, keyFieldTxnValue := range keyFieldTxnValue {
		presenceMap[keyFieldTxnValue] = make([]s.Presence, numOfClusters)
	}

	for index, cluster := range clusters {
		scoreResponse, err := cluster.Score(keyFieldTxnValue)
		if err != nil {
			log.Println("Repair Score Error", err)

			go w.Farm.instrumentation.RepairScoreError()
			continue
		}

		for keyMemeber, presence := range scoreResponse {
			presenceMap[keyMemeber][index] = presence
		}
	}

	var (
		inserts = map[int][]s.KeyFieldScoreTxnValue{}
		deletes = map[int][]s.KeyFieldScoreTxnValue{}
	)

	for keyFieldTxnValue, presenceSlice := range presenceMap {
		var (
			found        = false
			highestScore = float64(0)
			wasInserted  = false
		)

		for _, presence := range presenceSlice {
			if presence.Present && presence.Score >= highestScore {
				found = true
				highestScore = presence.Score
				wasInserted = wasInserted || presence.Inserted
			}
		}

		if !found {
			// This is need if the keyscore member has been removed by the time
			// the score has been asked for again
			continue
		}

		keyFieldScoreTxnValue := s.KeyFieldScoreTxnValue{
			Key:   keyFieldTxnValue.Key,
			Field: keyFieldTxnValue.Field,
			Score: highestScore,
			Value: keyFieldTxnValue.Value,
		}

		for index, presence := range presenceSlice {
			var (
				notThere = !presence.Present
				lowScore = presence.Score < highestScore
				wrongSet = presence.Inserted != wasInserted
			)

			if notThere || lowScore || wrongSet {
				if wasInserted {
					inserts[index] = append(inserts[index], keyFieldScoreTxnValue)
				} else {
					deletes[index] = append(deletes[index], keyFieldScoreTxnValue)
				}
			}
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(inserts) + len(deletes))

	errs := []string{}
	for index, keyFieldScoreTxnValues := range inserts {
		elements := clusters[index].Insert(keyFieldScoreTxnValues, maxSize)
		for e := range elements {
			func(e t.Element) {
				if err := t.ErrorFromElement(e); err != nil {
					errs = append(errs, err.Error())
				}
			}(e)
		}
	}

	for index, keyFieldScoreTxnValues := range deletes {
		elements := clusters[index].Delete(keyFieldScoreTxnValues, maxSize)
		for e := range elements {
			func(e t.Element) {
				defer wg.Done()

				if err := t.ErrorFromElement(e); err != nil {
					errs = append(errs, err.Error())
				}
			}(e)
		}
	}

	if timeout(wg, defaultTimeoutLatency) {
		go w.Farm.instrumentation.RepairError(1)
		return typex.Errorf(errors.Source, errors.Repair, "Repair Errors (Timeout)")
	}

	if len(errs) > 0 {
		go w.Farm.instrumentation.RepairError(len(errs))
		return typex.Errorf(errors.Source, errors.Repair, "Repair Errors (%s)", strings.Join(errs, ";"))
	}

	return nil
}

func timeout(wg *sync.WaitGroup, t time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return false
	case <-time.After(t):
		return true
	}
}
