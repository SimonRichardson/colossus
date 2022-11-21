package consul

import (
	"math/rand"
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/typex"
)

func Semaphore(s *Service, t Tactic) selectors.Semaphore {
	return semaphore{s, t}
}

type semaphore struct {
	service *Service
	tactic  Tactic
}

func (s semaphore) Lock(ns selectors.Namespace) (selectors.SemaphoreUnlock, error) {
	return s.write(func(c Cluster) <-chan sv.Element {
		return c.Lock(ns)
	})
}

func (s semaphore) write(fn func(Cluster) <-chan sv.Element) (selectors.SemaphoreUnlock, error) {
	var (
		service       = s.service
		clusters, err = selectClusters(service.clusters)

		retrieved = 0
		returned  = 0
	)
	if err != nil {
		return noopUnlock, err
	}

	began := beforeWrite(service.instrumentation, 1)
	defer afterWrite(service.instrumentation, began, retrieved, returned)

	var (
		elements = make(chan sv.Element, 1)
		errs     = []error{}
		changes  = []selectors.SemaphoreUnlock{}

		wg = &sync.WaitGroup{}
	)

	wg.Add(1)
	go func() { wg.Wait(); close(elements) }()

	scatterReads(s.tactic, clusters, fn, wg, elements)

	for element := range elements {
		retrieved++

		if err := sv.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}

		unlock := sv.SemaphoreUnlockFromElement(element)
		changes = append(changes, unlock)

		returned++
	}

	if len(errs) > 0 {
		return noopUnlock, typex.Errorf(errors.Source, errors.Complete,
			"Complete failure").With(errs...)
	}

	return head(changes)
}

func selectClusters(clusters []Cluster) ([]Cluster, error) {
	if len(clusters) < 1 {
		return nil, typex.Errorf(errors.Source, errors.UnexpectedArgument, "No clusters to locate.")
	}
	offset := rand.Intn(len(clusters))
	return clusters[offset : offset+1], nil
}

func head(x []selectors.SemaphoreUnlock) (selectors.SemaphoreUnlock, error) {
	if len(x) < 1 {
		return noopUnlock, typex.Errorf(errors.Source, errors.Complete,
			"Complete failure: no valid changes")
	}
	return x[0], nil
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

func beforeWrite(instr Instrumentation, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr.SemaphoreCall()
		instr.SemaphoreSendTo(numSends)
	}()
	return began
}

func afterWrite(instr Instrumentation, began time.Time, retrieved, returned int) {
	go func() {
		instr.SemaphoreDuration(time.Since(began))
		instr.SemaphoreRetrieved(retrieved)
		instr.SemaphoreReturned(returned)
	}()
}
