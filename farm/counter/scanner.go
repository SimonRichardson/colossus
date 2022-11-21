package counter

import (
	"sync"
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	t "github.com/SimonRichardson/colossus/cluster"
	r "github.com/SimonRichardson/colossus/cluster/counter"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/farm"
	"github.com/SimonRichardson/colossus/instrumentation"
	s "github.com/SimonRichardson/colossus/selectors"
)

// ScanAllReadAll defines a strategy to scan the store requesting to all the
// clusters and the reading back all the information from the clusters to
// perform a consensus
func ScanAllReadAll(f *Farm, t Tactic) s.Scanner {
	return scanAllReadAll{f, t}
}

type scanAllReadAll struct {
	*Farm
	tactic Tactic
}

func (w scanAllReadAll) Keys() ([]bs.Key, error) {
	return w.readKeys(func(c r.Cluster) <-chan t.Element {
		return c.Keys()
	})
}

func (w scanAllReadAll) Size(key bs.Key) (int, error) {
	return w.readInt(func(c r.Cluster) <-chan t.Element {
		return c.Size(key)
	})
}

func (w scanAllReadAll) Members(key bs.Key) ([]bs.Key, error) {
	return w.readKeys(func(c r.Cluster) <-chan t.Element {
		return c.Members(key)
	})
}

func (w scanAllReadAll) readKeys(fn func(r.Cluster) <-chan t.Element) ([]bs.Key, error) {
	var (
		clusters      = w.Farm.clusters
		numOfClusters = len(clusters)
	)

	began := beforeScan(w.Farm, numOfClusters)
	defer afterScan(w.Farm, began)

	var (
		elements = make(chan t.Element, numOfClusters)

		responses = []bs.Key{}
		retrieved = 0
		returned  = 0

		wg = sync.WaitGroup{}
	)

	wg.Add(numOfClusters)
	go func() { wg.Wait(); close(elements) }()

	if err := scatterReads(w.tactic, w.instrumentation, clusters, fn, &wg, elements); err != nil {
		return nil, err
	}

	for element := range elements {
		var (
			keys      = t.KeysFromElement(element)
			numOfKeys = len(keys)
		)

		retrieved += numOfKeys

		if err := t.ErrorFromElement(element); err != nil {
			return nil, err
		}

		if numOfKeys > 0 {
			returned += numOfKeys
			responses = append(responses, keys...)
		}
	}

	resultsScan(w.Farm, retrieved, returned)
	return responses, nil
}

func (w scanAllReadAll) readInt(fn func(r.Cluster) <-chan t.Element) (int, error) {
	var (
		clusters      = w.Farm.clusters
		numOfClusters = len(clusters)
	)

	began := beforeScan(w.Farm, numOfClusters)
	defer afterScan(w.Farm, began)

	var (
		elements = make(chan t.Element, numOfClusters)

		responses = []int{}
		retrieved = 0
		returned  = 0

		wg = sync.WaitGroup{}
	)

	wg.Add(numOfClusters)
	go func() { wg.Wait(); close(elements) }()

	if err := scatterReads(w.tactic, w.instrumentation, clusters, fn, &wg, elements); err != nil {
		return -1, err
	}

	for element := range elements {
		amount := t.AmountFromElement(element)
		retrieved++

		if err := t.ErrorFromElement(element); err != nil {
			return 0, err
		}

		returned++
		responses = append(responses, amount)
	}

	master := common.NewLargestInt()
	for _, v := range responses {
		master.Add(v)
	}

	defer resultsScan(w.Farm, retrieved, returned)

	response, err := master.Value()

	if err != nil {
		go w.Farm.instrumentation.ScanRepairNeeded(master.Len())
		return response, farm.NewPartialError(farm.Counter, response, responses)
	}

	return response, nil
}

func beforeScan(f *Farm, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr := f.instrumentation
		instr.ScanCall()
		instr.ScanSendTo(numSends)
	}()
	return began
}

func resultsScan(f *Farm, retrieved, returned int) {
	go func() {
		instr := f.instrumentation
		instr.ScanRetrieved(retrieved)
		instr.ScanReturned(returned)
	}()
}

func afterScan(f *Farm, began time.Time) {
	go func() {
		instr := f.instrumentation
		instr.ScanDuration(time.Since(began))
	}()
}

func scatterReads(
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
