package score

import (
	"sync"
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	t "github.com/SimonRichardson/colossus/cluster"
	r "github.com/SimonRichardson/colossus/colossus-shim/cluster/score"
	st "github.com/SimonRichardson/colossus/colossus-shim/selectors"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
)

// IncrementAllReadAll defines a strategy to write to all the cluster and then
// wait for all the cluster items to respond before continuing onwards.
func IncrementAllReadAll(f *Farm, t Tactic) st.Incrementer {
	return writeAllReadAll{f, t}
}

type writeAllReadAll struct {
	*Farm
	tactic Tactic
}

func (w writeAllReadAll) Increment(key bs.Key, time time.Time) (int, error) {
	return w.write(func(c r.Cluster) <-chan t.Element {
		return c.Increment(key)
	})
}

func (w writeAllReadAll) write(fn func(r.Cluster) <-chan t.Element) (int, error) {
	var (
		clusters      = w.Farm.clusters
		numOfClusters = len(clusters)

		retrieved = 0
		returned  = 0
	)

	began := beforeWrite(w.Farm, numOfClusters)
	defer afterWrite(w.Farm, began, retrieved, returned)

	var (
		elements = make(chan t.Element, numOfClusters)
		errs     = []error{}
		changes  = []int{}

		master = common.NewSimilarInt()
		repair = false

		wg = &sync.WaitGroup{}
	)

	wg.Add(numOfClusters)
	go func() { wg.Wait(); close(elements) }()

	scatterWrites(w.tactic, clusters, fn, wg, elements)

	for element := range elements {

		amount := t.AmountFromElement(element)
		retrieved += amount

		if err := t.ErrorFromElement(element); err != nil {
			repair = true
			errs = append(errs, err)
			continue
		}

		returned += amount

		// Detect if we need a read repair
		if !repair {
			repair = !master.Similar(amount)
		}
		changes = append(changes, amount)
	}

	if repair {
		// TODO : Log out error here!
	}
	if len(errs) > 0 {
		return -1, typex.Errorf(errors.Source, errors.Partial,
			"Partial Error (%s)", common.SumErrors(errs).Error())
	}

	return common.LargestInteger(changes), nil
}

func scatterWrites(
	tactic Tactic,
	clusters []r.Cluster,
	fn func(r.Cluster) <-chan t.Element,
	wg *sync.WaitGroup,
	dst chan t.Element,
) error {
	return tactic(clusters, func(c r.Cluster) {
		defer wg.Done()
		for e := range fn(c) {
			dst <- e
		}
	})
}

func beforeWrite(f *Farm, numSends int) time.Time {
	began := time.Now()
	go func() {
		/*
			// Do we care about this not happening :thinking_face:?
			instr := f.instrumentation
			instr.IncrementCall()
			instr.IncrementKeys(1)
			instr.IncrementSendTo(numSends)
		*/
	}()
	return began
}

func afterWrite(f *Farm, began time.Time, retrieved, returned int) {
	go func() {
		/*
			// Do we care about this not happening :thinking_face:?
			instr := f.instrumentation
			instr.IncrementRetrieved(retrieved)
			instr.IncrementReturned(returned)
			instr.IncrementDuration(time.Since(began))
		*/
	}()
}
