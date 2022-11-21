package store

import (
	"math"
	"sync"
	"time"

	t "github.com/SimonRichardson/colossus/cluster"
	r "github.com/SimonRichardson/colossus/cluster/store"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/farm"
	"github.com/SimonRichardson/colossus/instrumentation"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/typex"
)

// InsertAllReadAll defines a strategy to write to all the cluster and then
// wait for all the cluster items to respond before continuing onwards.
func InsertAllReadAll(f *Farm, t Tactic) s.Inserter {
	return writeAllReadAll{f, wtInsert, t}
}

// InsertAllReadQuorum defines a strategy to write to all the cluster and then
// wait for the quorum of cluster items to respond before continuing onwards.
func InsertAllReadQuorum(quorum float64) func(f *Farm, t Tactic) s.Inserter {
	return func(f *Farm, t Tactic) s.Inserter {
		return writeAllReadQuorum{f, wtInsert, t, quorum}
	}
}

// DeleteAllReadAll defines a strategy to write to all the cluster and then
// wait for all the cluster items to respond before continuing onwards.
func DeleteAllReadAll(f *Farm, t Tactic) s.Deleter {
	return writeAllReadAll{f, wtDelete, t}
}

// DeleteAllReadQuorum defines a strategy to write to all the cluster and then
// wait for the quorum of cluster items to respond before continuing onwards.
func DeleteAllReadQuorum(quorum float64) func(f *Farm, t Tactic) s.Deleter {
	return func(f *Farm, t Tactic) s.Deleter {
		return writeAllReadQuorum{f, wtDelete, t, quorum}
	}
}

type writeType int

const (
	wtInsert writeType = iota
	wtDelete
)

type writeAllReadAll struct {
	*Farm
	wtype  writeType
	tactic Tactic
}

func (w writeAllReadAll) Insert(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (int, error) {
	return w.write(func(c r.Cluster) <-chan t.Element {
		return c.Insert(members, maxSize)
	})
}

func (w writeAllReadAll) Delete(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (int, error) {
	return w.write(func(c r.Cluster) <-chan t.Element {
		return c.Delete(members, maxSize)
	})
}

func (w writeAllReadAll) Rollback(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) error {
	_, err := w.Delete(members, maxSize)
	return err
}

func (w writeAllReadAll) write(fn func(r.Cluster) <-chan t.Element) (int, error) {
	var (
		clusters      = w.Farm.clusters
		numOfClusters = len(clusters)

		retrieved = 0
		returned  = 0
	)

	began := beforeWrite(w.Farm, w.wtype, numOfClusters)
	defer afterWrite(w.Farm, w.wtype, began, retrieved, returned)

	var (
		elements = make(chan t.Element, numOfClusters)
		errs     = []error{}
		changes  = []int{}

		master  = common.NewSimilarInt()
		similar = true

		wg = &sync.WaitGroup{}
	)

	wg.Add(numOfClusters)
	go func() { wg.Wait(); close(elements) }()

	if err := scatterWrites(w.tactic, w.instrumentation, clusters, fn, wg, elements); err != nil {
		return -1, err
	}

	for element := range elements {
		amount := t.AmountFromElement(element)
		retrieved += amount

		if err := t.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}

		returned += amount

		// Detect if we need a read repair
		similar = similar && master.Similar(amount)
		changes = append(changes, amount)
	}

	// If the repair is fale, then go through it
	if len(errs) > 0 {
		repairWrite(w.Farm, w.wtype)

		return -1, typex.Errorf(errors.Source, errors.Partial,
			"Partial Error (%s)", common.SumErrors(errs).Error())
	} else if !similar {
		repairWrite(w.Farm, w.wtype)

		// We don't care about this error, we're just interested in the potential
		// value!
		value, err := master.Value()
		if err != nil {
			teleprinter.L.Error().Printf("Error from partial error %s\n", err)
		}
		return -1, farm.NewPartialError(farm.Counter, value, changes)
	}

	return common.AvgIntegers(changes), nil
}

type writeAllReadQuorum struct {
	*Farm
	wtype  writeType
	tactic Tactic
	quorum float64
}

func (w writeAllReadQuorum) Insert(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (int, error) {
	return w.write(func(c r.Cluster) <-chan t.Element {
		return c.Insert(members, maxSize)
	})
}

func (w writeAllReadQuorum) Delete(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) (int, error) {
	return w.write(func(c r.Cluster) <-chan t.Element {
		return c.Delete(members, maxSize)
	})
}

func (w writeAllReadQuorum) Rollback(members []s.KeyFieldScoreTxnValue, maxSize s.KeySizeExpiry) error {
	_, err := w.Delete(members, maxSize)
	return err
}

func (w writeAllReadQuorum) write(fn func(r.Cluster) <-chan t.Element) (int, error) {
	var (
		clusters           = w.Farm.clusters
		numOfClusters      = int(math.Ceil(float64(len(clusters)) * w.quorum))
		selected, deferred = selectClusters(clusters, numOfClusters)

		retrieved = 0
		returned  = 0
	)

	began := beforeWrite(w.Farm, w.wtype, numOfClusters)
	defer afterWrite(w.Farm, w.wtype, began, retrieved, returned)

	var (
		elements = make(chan t.Element, numOfClusters)
		errs     = []error{}
		changes  = []int{}

		master  = common.NewSimilarInt()
		similar = true

		wg = &sync.WaitGroup{}
	)

	wg.Add(numOfClusters)
	go func() { wg.Wait(); close(elements) }()

	if err := scatterWrites(w.tactic, w.instrumentation, selected, fn, wg, elements); err != nil {
		return -1, err
	}

	for element := range elements {
		amount := t.AmountFromElement(element)
		retrieved += amount

		if err := t.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}

		returned += amount

		// Detect if we need a read repair
		similar = similar && master.Similar(amount)
		changes = append(changes, amount)
	}

	// If the repair is fale, then go through it
	if len(errs) > 0 {
		repairWrite(w.Farm, w.wtype)

		return -1, typex.Errorf(errors.Source, errors.Partial,
			"Partial Error (%s)", common.SumErrors(errs).Error())
	} else if !similar {
		repairWrite(w.Farm, w.wtype)

		// We don't care about this error, we're just interested in the potential
		// value!
		value, err := master.Value()
		if err != nil {
			teleprinter.L.Error().Printf("Error from partial error %s\n", err)
		}
		return -1, farm.NewPartialError(farm.Store, value, changes)
	}

	// Not we only write to the rest of the quorum when we know that a repair
	// is not required.
	go func() {
		var (
			elements = make(chan t.Element, numOfClusters)
			errs     = []error{}
			wg       = &sync.WaitGroup{}
		)

		wg.Add(len(deferred))
		go func() { wg.Wait(); close(elements) }()

		if err := scatterWrites(w.tactic, w.instrumentation, deferred, fn, wg, elements); err != nil {
			return
		}

		for element := range elements {
			if err := t.ErrorFromElement(element); err != nil {
				errs = append(errs, err)
				continue
			}
		}

		if len(errs) > 0 {
			repairWrite(w.Farm, w.wtype)
		}
	}()

	return common.AvgIntegers(changes), nil
}

func scatterWrites(
	tactic Tactic,
	instr instrumentation.Instrumentation,
	clusters []r.Cluster,
	fn func(r.Cluster) <-chan t.Element,
	wg *sync.WaitGroup,
	dst chan t.Element,
) error {
	return tactic(clusters, func(k int, c r.Cluster) {
		began := time.Now()
		go instr.ClusterCall(k)
		defer func() {
			wg.Done()
			go instr.ClusterDuration(k, time.Since(began))
		}()

		for e := range fn(c) {
			dst <- e
		}
	})
}

func beforeWrite(f *Farm, wtype writeType, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr := f.instrumentation
		switch wtype {
		case wtInsert:
			instr.InsertCall()
			instr.InsertKeys(1)
			instr.InsertSendTo(numSends)
		case wtDelete:
			instr.DeleteCall()
			instr.DeleteKeys(1)
			instr.DeleteSendTo(numSends)
		}
	}()
	return began
}

func resultsWrite(f *Farm, wtype writeType) {
	go func() {
		instr := f.instrumentation
		switch wtype {
		case wtInsert:
			instr.InsertQuorumFailure()
		case wtDelete:
			instr.DeleteQuorumFailure()
		}
	}()
}

func repairWrite(f *Farm, wtype writeType) {
	go func() {
		instr := f.instrumentation
		switch wtype {
		case wtInsert:
			instr.InsertRepairRequired()
		case wtDelete:
			instr.DeleteRepairRequired()
		}
	}()
}

func afterWrite(f *Farm, wtype writeType, began time.Time, retrieved, returned int) {
	go func() {
		instr := f.instrumentation
		switch wtype {
		case wtInsert:
			instr.InsertRetrieved(retrieved)
			instr.InsertReturned(returned)
			instr.InsertDuration(time.Since(began))
		case wtDelete:
			instr.DeleteRetrieved(retrieved)
			instr.DeleteReturned(returned)
			instr.DeleteDuration(time.Since(began))
		}
	}()
}
