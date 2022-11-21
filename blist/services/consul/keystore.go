package consul

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	fs "github.com/SimonRichardson/colossus/fusion/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

func KeyStore(s *Service, t Tactic) selectors.KeyStore {
	return keyStore{s, t}
}

type keyStore struct {
	service *Service
	tactic  Tactic
}

func (s keyStore) List(p fs.Prefix) (map[string]int, error) {
	return s.write(func(c Cluster) <-chan sv.Element {
		return c.List(p)
	})
}

func (s keyStore) write(fn func(Cluster) <-chan sv.Element) (map[string]int, error) {
	var (
		service       = s.service
		clusters, err = selectClusters(service.clusters)

		retrieved = 0
		returned  = 0
	)
	if err != nil {
		return nil, err
	}

	began := beforeWrite(service.instrumentation, 1)
	defer afterWrite(service.instrumentation, began, retrieved, returned)

	var (
		elements = make(chan sv.Element, 1)
		errs     = []error{}
		changes  = map[string]int{}

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

		values := sv.MapStringIntFromElement(element)
		for k, v := range values {
			changes[k] = v
		}

		returned += len(values)
	}

	if len(errs) > 0 {
		return nil, typex.Errorf(errors.Source, errors.Complete,
			"Complete failure").With(errs...)
	}

	return changes, nil
}

func beforeKeyStoreWrite(instr Instrumentation, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr.KeyStoreCall()
		instr.KeyStoreSendTo(numSends)
	}()
	return began
}

func afterKeyStoreWrite(instr Instrumentation, began time.Time, retrieved, returned int) {
	go func() {
		instr.KeyStoreDuration(time.Since(began))
		instr.KeyStoreRetrieved(retrieved)
		instr.KeyStoreReturned(returned)
	}()
}
