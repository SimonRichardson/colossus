package store

import (
	"math"
	"math/rand"
	"sync"
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	t "github.com/SimonRichardson/colossus/cluster"
	r "github.com/SimonRichardson/colossus/cluster/store"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/farm"
	"github.com/SimonRichardson/colossus/instrumentation"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

const (
	defaultTimeoutLatency = time.Second

	defaultLimit   = 10
	defaultMaxSize = 10
	defaultExpiry  = time.Minute
)

// SelectOneReadOne defines a strategy to write to all the cluster and then
// wait for all the cluster items to respond before continuing onwards.
func SelectOneReadOne(f *Farm, t Tactic) s.Selector {
	return selectOneReadOne{f, t}
}

// SelectQuorumReadAll defines a strategy to write to just the minimum quorum of
// servers and then read all the results back.
func SelectQuorumReadAll(quorum float64) func(f *Farm, t Tactic) s.Selector {
	return func(f *Farm, t Tactic) s.Selector {
		return selectQuorumReadAll{f, t, quorum}
	}
}

// SelectAllReadAll defines a strategy to write to all of servers and then read
// all the results back.
func SelectAllReadAll(f *Farm, t Tactic) s.Selector {
	return selectQuorumReadAll{f, t, 1.0}
}

type selectOneReadOne struct {
	*Farm
	tactic Tactic
}

func (w selectOneReadOne) Select(key bs.Key, field bs.Key) (s.KeyFieldScoreTxnValue, error) {
	return unwrapSelection(w.read(key, func(c r.Cluster) <-chan t.Element {
		return c.Select(key, field)
	}))
}

func (w selectOneReadOne) SelectRange(key bs.Key, limit int, maxSize s.KeySizeExpiry) ([]s.KeyFieldScoreTxnValue, error) {
	return w.read(key, func(c r.Cluster) <-chan t.Element {
		return c.SelectRange(key, limit, maxSize)
	})
}

func (w selectOneReadOne) read(key bs.Key,
	fn func(r.Cluster) <-chan t.Element,
) ([]s.KeyFieldScoreTxnValue, error) {
	var (
		clusters      = w.Farm.clusters
		numOfClusters = 1
	)

	began := beforeRead(w.Farm, 1, numOfClusters)
	defer afterRead(w.Farm, began)

	var (
		selected = []r.Cluster{clusters[rand.Intn(len(clusters))]}
		elements = send(key, w.tactic, w.instrumentation, selected, numOfClusters, fn)

		response  = []s.KeyFieldScoreTxnValue{}
		retrieved = 0
		returned  = 0
	)

	for element := range elements {
		var (
			members      = t.ValuesFromElement(element)
			numOfMembers = len(members)
		)

		retrieved += numOfMembers

		if err := t.ErrorFromElement(element); err != nil {
			return nil, err
		}

		if numOfMembers > 0 {
			returned += numOfMembers
			response = append(response, members...)
		}
		break
	}

	resultsRead(w.Farm, retrieved, returned)
	return response, nil
}

type selectQuorumReadAll struct {
	*Farm
	tactic Tactic
	quorum float64
}

func (w selectQuorumReadAll) Select(key bs.Key, field bs.Key) (s.KeyFieldScoreTxnValue, error) {
	return unwrapSelection(w.read(func(c r.Cluster) <-chan t.Element {
		return c.Select(key, field)
	}, defaultLimit, s.MakeKeySizeSingleton(key, defaultMaxSize, defaultExpiry)))
}

func (w selectQuorumReadAll) SelectRange(key bs.Key, limit int, maxSize s.KeySizeExpiry) ([]s.KeyFieldScoreTxnValue, error) {
	return w.read(func(c r.Cluster) <-chan t.Element {
		return c.SelectRange(key, limit, maxSize)
	}, limit, maxSize)
}

func (w selectQuorumReadAll) read(fn func(r.Cluster) <-chan t.Element, limit int, maxSize s.KeySizeExpiry) ([]s.KeyFieldScoreTxnValue, error) {
	var (
		clusters      = w.Farm.clusters
		numOfClusters = int(math.Ceil(float64(len(clusters)) * w.quorum))
	)

	began := beforeRead(w.Farm, 1, numOfClusters)
	defer afterRead(w.Farm, began)

	var (
		elements    = make(chan t.Element, numOfClusters)
		selected, _ = selectClusters(clusters, numOfClusters)

		responses = []farm.TupleSet{}
		retrieved = 0
		returned  = 0

		wg = &sync.WaitGroup{}
	)

	wg.Add(numOfClusters)

	go func() { wg.Wait(); close(elements) }()

	if err := scatterReads(w.tactic, w.instrumentation, selected, fn, wg, elements); err != nil {
		return []s.KeyFieldScoreTxnValue{}, err
	}

	for element := range elements {
		var (
			members      = t.ValuesFromElement(element)
			numOfMembers = len(members)
		)

		retrieved += numOfMembers

		if err := t.ErrorFromElement(element); err != nil {
			return nil, err
		}

		if numOfMembers > 0 {
			returned += numOfMembers
			responses = append(responses, farm.MakeSet(members))
		}
	}

	var (
		union, difference = farm.UnionDifference(responses)
		response          = union.OrderedLimitedSlice(limit)

		repairs = farm.KeyFieldTxnValueSet{}
	)

	repairs.AddMany(difference)

	if len(repairs) > 0 {
		go w.Farm.Repair(repairs.Slice(), maxSize)
	}

	resultsRead(w.Farm, retrieved, returned)

	return response, nil
}

func selectClusters(clusters []r.Cluster, split int) ([]r.Cluster, []r.Cluster) {
	num := len(clusters)
	if split >= num {
		return clusters, []r.Cluster{}
	}

	var (
		dest = make([]r.Cluster, num)
		perm = rand.Perm(num)
	)
	for i, v := range perm {
		dest[v] = clusters[i]
	}

	return dest[:split], dest[split:]
}

func beforeRead(f *Farm, numKeys, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr := f.instrumentation
		instr.SelectCall()
		instr.SelectKeys(numKeys)
		instr.SelectSendTo(numSends)
	}()
	return began
}

func resultsRead(f *Farm, retrieved, returned int) {
	go func() {
		instr := f.instrumentation
		instr.SelectRetrieved(retrieved)
		instr.SelectReturned(returned)
	}()
}

func afterRead(f *Farm, began time.Time) {
	go func() {
		instr := f.instrumentation
		instr.SelectDuration(time.Since(began))
	}()
}

func send(key bs.Key,
	tactic Tactic,
	instr instrumentation.Instrumentation,
	clusters []r.Cluster,
	waitFor int,
	fn func(r.Cluster) <-chan t.Element,
) <-chan t.Element {
	elements := make(chan t.Element, waitFor)

	wg := sync.WaitGroup{}
	wg.Add(waitFor)
	go func() { wg.Wait(); close(elements) }()

	if err := scatterReads(tactic, instr, clusters, fn, &wg, elements); err != nil {
		elements <- t.NewErrorElement(key, err)
	}

	return elements
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

func unwrapSelection(sel []s.KeyFieldScoreTxnValue, err error) (s.KeyFieldScoreTxnValue, error) {
	if err != nil {
		return s.KeyFieldScoreTxnValue{}, err
	}
	if len(sel) < 1 {
		return s.KeyFieldScoreTxnValue{}, typex.Errorf(errors.Source, errors.UnexpectedResults, "Not found")
	}
	return sel[0], nil
}
