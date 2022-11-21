package cache

import (
	"math/rand"
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/typex"
)

func Encoder(s *Service, t Tactic) selectors.Encoder {
	return encoder{s, t}
}

type encoder struct {
	*Service
	tactic Tactic
}

func (e encoder) GetBytes(key bs.Key) ([]byte, error) {
	return e.read(func(c Cluster) <-chan sv.Element {
		return c.GetBytes(key)
	})
}

func (e encoder) SetBytes(key bs.Key, bytes []byte) error {
	return e.write(func(c Cluster) <-chan sv.Element {
		return c.SetBytes(key, bytes)
	})
}

func (e encoder) DelBytes(key bs.Key) error {
	return e.write(func(c Cluster) <-chan sv.Element {
		return c.DelBytes(key)
	})
}

func (e encoder) read(fn func(Cluster) <-chan sv.Element) ([]byte, error) {
	cluster, err := require(e.clusters)
	if err != nil {
		return nil, err
	}

	began := beforeWrite(e.instrumentation, 1)
	defer afterWrite(e.instrumentation, began)

	var (
		elements = make(chan sv.Element, 1)
		errs     = make([]error, 0)
		changes  = make([][]byte, 0)

		wg = &sync.WaitGroup{}
	)

	wg.Add(1)
	go func() { wg.Wait(); close(elements) }()

	scatterReads(e.tactic, []Cluster{cluster}, fn, wg, elements)

	for element := range elements {
		if err := sv.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}

		event := sv.BytesFromElement(element)
		changes = append(changes, event)
	}

	if len(errs) > 0 {
		return nil, typex.Errorf(errors.Source, errors.Complete,
			"Complete failure").With(errs...)
	}
	return headBytes(changes)
}

func (e encoder) write(fn func(Cluster) <-chan sv.Element) error {
	numClusters := len(e.clusters)

	began := beforeWrite(e.instrumentation, numClusters)
	defer afterWrite(e.instrumentation, began)

	var (
		elements = make(chan sv.Element, numClusters)
		errs     = make([]error, 0)

		wg = &sync.WaitGroup{}
	)

	wg.Add(numClusters)
	go func() { wg.Wait(); close(elements) }()

	scatterReads(e.tactic, e.clusters, fn, wg, elements)

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

func require(clusters []Cluster) (Cluster, error) {
	if len(clusters) < 1 {
		return nil, typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"No cluster to locate.")
	}
	return clusters[rand.Intn(len(clusters))], nil
}

func headBytes(bytes [][]byte) ([]byte, error) {
	if len(bytes) < 1 {
		return nil, typex.Errorf(errors.Source, errors.UnexpectedResults,
			"Unexpected number of items.")
	}
	return bytes[0], nil
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
		instr.EncodeCall()
		instr.EncodeSendTo(numSends)
	}()
	return began
}

func afterWrite(instr Instrumentation, began time.Time) {
	go func() {
		instr.EncodeDuration(time.Since(began))
	}()
}
