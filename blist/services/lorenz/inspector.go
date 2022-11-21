package lorenz

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/typex"
)

// Inspector defines a strategy to read all the versions from the clusters
func Inspector(s *Service, t Tactic) selectors.Inspector {
	return inspector{s, t}
}

type inspector struct {
	*Service
	tactic Tactic
}

func (s inspector) Version() (map[string][]selectors.Version, error) {
	return s.readVersion(func(c Cluster) <-chan sv.Element {
		return c.Version()
	})
}

func (s inspector) readVersion(fn func(Cluster) <-chan sv.Element) (map[string][]selectors.Version, error) {
	var (
		clusters      = s.clusters
		numOfClusters = len(clusters)

		retrieved = 0
		returned  = 0
	)

	began := beforeInspectRead(s.instrumentation, numOfClusters)
	defer afterInspectRead(s.instrumentation, began, retrieved, returned)

	var (
		elements = make(chan sv.Element, numOfClusters)
		errs     = []error{}
		changes  = []selectors.Version{}

		wg = &sync.WaitGroup{}
	)

	wg.Add(numOfClusters)
	go func() { wg.Wait(); close(elements) }()

	scatterReads(s.tactic, clusters, fn, wg, elements)

	for element := range elements {
		retrieved++

		if err := sv.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}

		returned++

		version := sv.VersionFromElement(element)
		changes = append(changes, version)
	}

	if len(errs) > 0 {
		return nil, typex.Errorf(errors.Source, errors.Complete,
			"Complete failure").With(errs...)
	}

	return map[string][]selectors.Version{
		defaultServiceName: changes,
	}, nil
}

func beforeInspectRead(instr Instrumentation, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr.InspectCall()
		instr.InspectSendTo(numSends)
	}()
	return began
}

func afterInspectRead(instr Instrumentation, began time.Time, retrieved, returned int) {
	go func() {
		instr.InspectDuration(time.Since(began))
		instr.InspectRetrieved(retrieved)
		instr.InspectReturned(returned)
	}()
}
