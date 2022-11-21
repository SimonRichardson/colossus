package lorenz

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/typex"
)

func EventSelector(s *Service, t Tactic) selectors.EventSelector {
	return selector{s, t, ioEvent, &noopCache{}}
}

func EventCacheSelector(s *Service, t Tactic) selectors.EventSelector {
	return selector{s, t, ioEvent, newEventCache()}
}

func CodeSetSelector(s *Service, t Tactic) selectors.CodeSetSelector {
	return selector{s, t, ioCodeSet, &noopCache{}}
}

type ioType int

const (
	ioEvent ioType = iota
	ioCodeSet
)

type selector struct {
	service    *Service
	tactic     Tactic
	ioType     ioType
	eventCache EventCache
}

func (s selector) SelectEventByKey(key selectors.Key) (selectors.Event, error) {
	return s.readEvent(key, func(c Cluster) <-chan sv.Element {
		return c.SelectEventByKey(key)
	})
}

func (s selector) SelectEventsByOffset(offset, limit int) ([]selectors.Event, error) {
	return s.readEvents(func(c Cluster) <-chan sv.Element {
		return c.SelectEventsByOffset(offset, limit)
	})
}

func (s selector) SelectCodeForEvent(event selectors.Event, user selectors.User) (selectors.CodeSet, error) {
	return s.readCode(func(c Cluster) <-chan sv.Element {
		return c.SelectCodeForEvent(event, user)
	})
}

func (s selector) readEvents(fn func(Cluster) <-chan sv.Element) ([]selectors.Event, error) {
	var (
		clusters      = s.service.clusters
		numOfClusters = len(s.service.clusters)
	)

	began := beforeRead(s.service.instrumentation, s.ioType, 1)
	defer afterRead(s.service.instrumentation, s.ioType, began)

	var (
		elements = make(chan sv.Element)
		errs     = make([]error, 0)
		changes  = make([]selectors.Event, 0)

		wg = &sync.WaitGroup{}
	)

	wg.Add(numOfClusters)
	go func() { wg.Wait(); close(elements) }()

	go func() {
		if err := s.tactic(clusters, func(n Cluster) {
			defer wg.Done()
			for e := range fn(n) {
				elements <- e
			}
		}); err != nil {
			elements <- sv.NewErrorElement(err)
			wg.Done()
		}
	}()

	for element := range elements {
		if err := sv.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}

		value := sv.EventsFromElement(element)
		changes = append(changes, value...)
	}

	if len(errs) > 0 {
		return nil, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Error selecting Event")
	}

	return changes, nil
}

func (s selector) readEvent(key selectors.Key, fn func(Cluster) <-chan sv.Element) (selectors.Event, error) {
	var (
		clusters      = s.service.clusters
		numOfClusters = len(s.service.clusters)
	)

	began := beforeRead(s.service.instrumentation, s.ioType, 1)
	defer afterRead(s.service.instrumentation, s.ioType, began)

	// Exit out early!
	if usr, ok := s.eventCache.Get(key); ok {
		return usr, nil
	}

	var (
		elements = make(chan sv.Element)
		errs     = make([]error, 0)
		changes  = make([]selectors.Event, 0)

		wg = &sync.WaitGroup{}
	)

	wg.Add(numOfClusters)
	go func() { wg.Wait(); close(elements) }()

	go func() {
		if err := s.tactic(clusters, func(n Cluster) {
			defer wg.Done()
			for e := range fn(n) {
				elements <- e
			}
		}); err != nil {
			elements <- sv.NewErrorElement(err)
			wg.Done()
		}
	}()

	// This is a race!
	for element := range elements {
		if err := sv.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}

		value := sv.EventFromElement(element)
		changes = append(changes, value)
	}

	if len(errs) > 0 {
		return selectors.Event{}, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Error selecting Event")
	}

	if len(changes) < 1 {
		return selectors.Event{}, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Event not found from key.")
	}

	event, err := eventsHead(changes)
	if err == nil {
		// Save the event to the cache
		go s.eventCache.Set(key, event, defaultEventCacheDuration)
	}
	return event, err
}

func (s selector) readCode(fn func(Cluster) <-chan sv.Element) (selectors.CodeSet, error) {
	var (
		clusters      = s.service.clusters
		numOfClusters = len(s.service.clusters)
	)

	began := beforeRead(s.service.instrumentation, s.ioType, numOfClusters)
	defer afterRead(s.service.instrumentation, s.ioType, began)

	var (
		elements = make(chan sv.Element)
		errs     = make([]error, 0)
		changes  = make([]selectors.CodeSet, 0)

		wg = &sync.WaitGroup{}
	)

	wg.Add(numOfClusters)
	go func() { wg.Wait(); close(elements) }()

	go func() {
		if err := s.tactic(clusters, func(n Cluster) {
			defer wg.Done()
			for e := range fn(n) {
				elements <- e
			}
		}); err != nil {
			elements <- sv.NewErrorElement(err)
			wg.Done()
		}
	}()

	// This is a race!
	for element := range elements {
		if err := sv.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}

		value := sv.CodeSetFromElement(element)
		changes = append(changes, value)

		break
	}

	if len(errs) > 0 {
		return selectors.CodeSet{}, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Error selecting CodeSet")
	}

	if len(changes) < 1 {
		return selectors.CodeSet{}, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"CodeSet not found from key.")
	}

	return codesHead(changes)
}

func eventsHead(changes []selectors.Event) (selectors.Event, error) {
	if len(changes) < 1 {
		return selectors.Event{}, typex.Errorf(errors.Source, errors.NoCaseFound,
			"Event not found")
	}
	return changes[0], nil
}

func codesHead(changes []selectors.CodeSet) (selectors.CodeSet, error) {
	if len(changes) < 1 {
		return selectors.CodeSet{}, typex.Errorf(errors.Source, errors.NoCaseFound,
			"CodeSet not found")
	}
	return changes[0], nil
}

func beforeRead(instr Instrumentation, t ioType, numSends int) time.Time {
	began := time.Now()
	go func() {
		switch t {
		case ioEvent:
			instr.EventSelectCall()
			instr.EventSelectSendTo(numSends)
		case ioCodeSet:
			instr.CodeSelectCall()
			instr.CodeSelectSendTo(numSends)
		}

	}()
	return began
}

func afterRead(instr Instrumentation, t ioType, began time.Time) {
	go func() {
		switch t {
		case ioEvent:
			instr.EventSelectDuration(time.Since(began))
		case ioCodeSet:
			instr.CodeSelectDuration(time.Since(began))
		}
	}()
}
