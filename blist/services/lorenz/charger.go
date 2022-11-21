package lorenz

import (
	"math/rand"
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/typex"
)

// Charger represents a strategy to release all the items from the clusters
func Charger(s *Service, t Tactic) selectors.Charger {
	return charger{s, t}
}

type charger struct {
	*Service
	tactic Tactic
}

func (s charger) Charge(event selectors.Event,
	user selectors.User,
	element selectors.Payment,
) (selectors.Key, error) {
	return s.write(func(c Cluster) <-chan sv.Element {
		return c.Charge(event, user, element)
	})
}

func (s charger) write(fn func(Cluster) <-chan sv.Element) (selectors.Key, error) {
	var (
		clusters, err = selectClusters(s.clusters)

		retrieved = 0
		returned  = 0
	)
	if err != nil {
		return selectors.Key(""), err
	}

	began := beforeChargeWrite(s.instrumentation, 1)
	defer afterChargeWrite(s.instrumentation, began, retrieved, returned)

	var (
		numOfClusters = len(clusters)

		elements = make(chan sv.Element, numOfClusters)
		errs     = []error{}
		changes  = []selectors.Key{}

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

		key := sv.KeyFromElement(element)
		changes = append(changes, key)

		returned++
	}

	if len(errs) > 0 {
		return selectors.Key(""), typex.Errorf(errors.Source, errors.Complete,
			"Complete failure").With(errs...)
	}

	return head(changes)
}

func selectClusters(clusters []Cluster) ([]Cluster, error) {
	if len(clusters) < 1 {
		return nil, typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"No cluster to locate.")
	}
	offset := rand.Intn(len(clusters))
	return clusters[offset : offset+1], nil
}

func head(x []selectors.Key) (selectors.Key, error) {
	if len(x) < 1 {
		return selectors.Key(""), typex.Errorf(errors.Source, errors.Complete,
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

func beforeChargeWrite(instr Instrumentation, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr.ChargeCall()
		instr.ChargeSendTo(numSends)
	}()
	return began
}

func afterChargeWrite(instr Instrumentation, began time.Time, retrieved, returned int) {
	go func() {
		instr.ChargeDuration(time.Since(began))
		instr.ChargeRetrieved(retrieved)
		instr.ChargeReturned(returned)
	}()
}
