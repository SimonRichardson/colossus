package resque

import (
	"math/rand"
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/typex"
)

type enqType int

const (
	etEnq enqType = iota
	etReg
)

func Enqueuer(s *Service, t Tactic) selectors.Enqueuer {
	return enqueuer{s, t, etEnq}
}

func Register(s *Service, t Tactic) selectors.Register {
	return enqueuer{s, t, etReg}
}

type enqueuer struct {
	*Service
	tactic  Tactic
	enqType enqType
}

func (e enqueuer) EnqueueBytes(queue selectors.Queue, class selectors.Class, value []byte) error {
	return e.write(func(c Cluster) <-chan sv.Element {
		return c.EnqueueBytes(queue, class, value)
	})
}

func (e enqueuer) DequeueBytes(queue selectors.Queue, class selectors.Class) ([]byte, error) {
	return e.read(func(c Cluster) <-chan sv.Element {
		return c.DequeueBytes(queue, class)
	})
}

func (e enqueuer) RegisterFailure(queue selectors.Queue,
	class selectors.Class,
	failure selectors.Failure,
) error {
	return e.write(func(c Cluster) <-chan sv.Element {
		return c.RegisterFailure(queue, class, failure)
	})
}

func (e enqueuer) read(fn func(Cluster) <-chan sv.Element) ([]byte, error) {
	cluster, err := require(e.clusters)
	if err != nil {
		return nil, err
	}

	began := beforeRead(e.instrumentation, 1)
	defer afterRead(e.instrumentation, began)

	var (
		elements = make(chan sv.Element, 1)
		errs     = []error{}
		changes  = [][]byte{}

		wg = &sync.WaitGroup{}
	)

	wg.Add(1)
	go func() { wg.Wait(); close(elements) }()

	if err := scatterReads(e.tactic, []Cluster{cluster}, fn, wg, elements); err != nil {
		return nil, typex.Errorf(errors.Source, errors.Complete,
			"Read failure").With(err)
	}

	for element := range elements {
		if err := sv.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}

		value := sv.BytesFromElement(element)
		changes = append(changes, value)
	}

	if len(errs) > 0 {
		return nil, typex.Errorf(errors.Source, errors.Complete,
			"Complete failure").With(errs...)
	}
	return headBytes(changes)
}

func (e enqueuer) write(fn func(Cluster) <-chan sv.Element) error {
	var (
		clusters    = e.clusters
		numClusters = len(e.clusters)
	)

	began := beforeWrite(e.enqType, e.instrumentation, numClusters)
	defer afterWrite(e.enqType, e.instrumentation, began)

	var (
		elements = make(chan sv.Element, numClusters)
		errs     = make([]error, 0)

		wg = &sync.WaitGroup{}
	)

	wg.Add(numClusters)
	go func() { wg.Wait(); close(elements) }()

	if err := scatterReads(e.tactic, clusters, fn, wg, elements); err != nil {
		return typex.Errorf(errors.Source, errors.Complete,
			"Write failure").With(err)
	}

	for element := range elements {
		if err := sv.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}
	}

	if len(errs) > 0 {
		return typex.Errorf(errors.Source, errors.Complete,
			"Complete failure").With(errs...)
	}
	return nil
}

func headBytes(values [][]byte) ([]byte, error) {
	if len(values) < 1 {
		return nil, typex.Errorf(errors.Source, errors.NoCaseFound,
			"No value found.")
	}
	return values[0], nil
}

func require(clusters []Cluster) (Cluster, error) {
	if len(clusters) < 1 {
		return nil, typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"No cluster to locate.")
	}
	return clusters[rand.Intn(len(clusters))], nil
}

func scatterReads(tactic Tactic,
	clusters []Cluster,
	fn func(Cluster) <-chan sv.Element,
	wg *sync.WaitGroup,
	dst chan sv.Element,
) error {
	return tactic(clusters, func(n Cluster) {
		defer wg.Done()
		for e := range fn(n) {
			dst <- e
		}
	})
}

func beforeWrite(t enqType, instr Instrumentation, numSends int) time.Time {
	began := time.Now()
	go func() {
		switch t {
		case etEnq:
			instr.EnqueueCall()
			instr.EnqueueSendTo(numSends)
		case etReg:
			instr.RegisterFailureCall()
			instr.RegisterFailureSendTo(numSends)
		}
	}()
	return began
}

func afterWrite(t enqType, instr Instrumentation, began time.Time) {
	go func() {
		switch t {
		case etEnq:
			instr.EnqueueDuration(time.Since(began))
		case etReg:
			instr.RegisterFailureDuration(time.Since(began))
		}
	}()
}

func beforeRead(instr Instrumentation, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr.DequeueCall()
		instr.DequeueSendTo(numSends)
	}()
	return began
}

func afterRead(instr Instrumentation, began time.Time) {
	go func() {
		instr.DequeueDuration(time.Since(began))
	}()
}
