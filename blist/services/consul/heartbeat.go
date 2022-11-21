package consul

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/typex"
)

func Heartbeat(s *Service, t Tactic) selectors.Heartbeat {
	return heartbeat{s, t}
}

type heartbeat struct {
	service *Service
	tactic  Tactic
}

func (s heartbeat) Heartbeat(ns selectors.HealthStatus) error {
	return s.write(func(c Cluster) <-chan sv.Element {
		return c.Heartbeat(ns)
	})
}

func (s heartbeat) write(fn func(Cluster) <-chan sv.Element) error {
	var (
		service       = s.service
		clusters, err = selectClusters(service.clusters)

		retrieved = 0
		returned  = 0
	)
	if err != nil {
		return err
	}

	began := beforeWrite(service.instrumentation, 1)
	defer afterWrite(service.instrumentation, began, retrieved, returned)

	var (
		elements = make(chan sv.Element, 1)
		errs     = []error{}

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

		returned++
	}

	if len(errs) > 0 {
		return typex.Errorf(errors.Source, errors.Complete,
			"Complete failure").With(errs...)
	}

	return nil
}

func beforeHeartbeatWrite(instr Instrumentation, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr.HeartbeatCall()
		instr.HeartbeatSendTo(numSends)
	}()
	return began
}

func afterHeartbeatWrite(instr Instrumentation, began time.Time, retrieved, returned int) {
	go func() {
		instr.HeartbeatDuration(time.Since(began))
		instr.HeartbeatRetrieved(retrieved)
		instr.HeartbeatReturned(returned)
	}()
}
