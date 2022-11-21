package notifier

import (
	"log"
	"math/rand"
	"sync"
	"time"

	t "github.com/SimonRichardson/colossus/cluster"
	r "github.com/SimonRichardson/colossus/cluster/notifier"
	"github.com/SimonRichardson/colossus/common"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/semaphore"
)

const (
	defaultSize    = 5000
	defaultTimeout = time.Second * 10
)

// Individual defines a strategy to write to all the cluster and then
// wait for all the cluster items to respond before continuing onwards.
func Individual(f *Farm, t Tactic) s.Notifier {
	return individual{f, t}
}

// Bulk defines a strategy to write to all the cluster and then wait for all the
// cluster items to respond before continuing onwards.
func Bulk(f *Farm, t Tactic) s.Notifier {
	var (
		individual = Individual(f, t)
		sem        = semaphore.NewBulk(func(values []semaphore.BulkItem) {
			m := toBulkItems(values).bucketize()
			for _, v := range m {
				for _, x := range v {
					individual.Publish(x.channel, x.members)
				}
			}
		}, defaultSize, defaultTimeout)
	)
	return bulk{f, sem, individual}
}

type individual struct {
	*Farm
	tactic Tactic
}

func (w individual) Publish(channel s.Channel, members []s.KeyFieldScoreSizeExpiry) error {
	return w.write(func(c r.Cluster) <-chan t.Element {
		return c.Publish(channel, members)
	})
}

func (w individual) Unpublish(channel s.Channel, members []s.KeyFieldScoreSizeExpiry) error {
	return w.write(func(c r.Cluster) <-chan t.Element {
		return c.Unpublish(channel, members)
	})
}

func (w individual) Subscribe(channel s.Channel) <-chan s.KeyFieldScoreSizeExpiry {
	return w.read(func(c r.Cluster) <-chan t.Element {
		return c.Subscribe(channel)
	})
}

func (w individual) write(fn func(r.Cluster) <-chan t.Element) error {
	var (
		clusters      = selectClusters(w.Farm.clusters)
		numOfClusters = len(clusters)
	)

	began := before(w.Farm, numOfClusters)
	defer after(w.Farm, began)

	var (
		elements  = make(chan t.Element, numOfClusters)
		errors    = []error{}
		retrieved = 0
		returned  = 0

		wg = sync.WaitGroup{}
	)

	wg.Add(numOfClusters)
	go func() { wg.Wait(); close(elements) }()

	// distribute randomly across the cluster
	scatterWrites(w.tactic, clusters, fn, &wg, elements)

	for element := range elements {
		retrieved++

		if err := t.ErrorFromElement(element); err != nil {
			errors = append(errors, err)
			continue
		}

		returned++
	}

	results(w.Farm, retrieved, returned)

	if len(errors) > 0 {
		return common.SumErrors(errors)
	}
	return nil
}

func (w individual) read(fn func(r.Cluster) <-chan t.Element) <-chan s.KeyFieldScoreSizeExpiry {
	var (
		clusters      = w.Farm.clusters
		numOfClusters = len(clusters)

		out = make(chan s.KeyFieldScoreSizeExpiry)
	)

	began := before(w.Farm, numOfClusters)
	defer after(w.Farm, began)

	for _, v := range clusters {

		go func(cluster r.Cluster) {
			for element := range fn(cluster) {
				if err := t.ErrorFromElement(element); err != nil {
					return
				}

				keyFieldScoreSizeExpiry := t.KeyFieldScoreSizeExpiryFromElement(element)
				if keyFieldScoreSizeExpiry.Key.Len() < 1 {
					continue
				}

				out <- keyFieldScoreSizeExpiry
			}
		}(v)
	}

	return out
}

type bulk struct {
	*Farm
	bulk       *semaphore.Bulk
	individual s.Notifier
}

func (w bulk) Publish(channel s.Channel, members []s.KeyFieldScoreSizeExpiry) error {
	return w.bulk.Add(bulkItem{channel, members})
}

func (w bulk) Unpublish(channel s.Channel, members []s.KeyFieldScoreSizeExpiry) error {
	go w.bulk.Remove(bulkItem{channel, members})
	return w.individual.Unpublish(channel, members)
}

func (w bulk) Subscribe(channel s.Channel) <-chan s.KeyFieldScoreSizeExpiry {
	return w.individual.Subscribe(channel)
}

type bulkItem struct {
	channel s.Channel
	members []s.KeyFieldScoreSizeExpiry
}

func (i bulkItem) Len() int {
	return len(i.members)
}

func (i bulkItem) Key() string {
	return i.channel.String()
}

func (i bulkItem) Equals(other semaphore.BulkItem) bool {
	// Try and do best guess first!
	if x, ok := other.(bulkItem); ok && i.channel == x.channel && i.Len() == x.Len() {
		for k, v := range i.members {
			z := x.members[k]
			if v.Key != z.Key || v.Field != z.Field || v.Score != z.Score {
				return false
			}
		}
		return true
	}
	return false
}

type bulkItems []bulkItem

func toBulkItems(values []semaphore.BulkItem) bulkItems {
	res := make([]bulkItem, 0, len(values))
	for _, v := range values {
		if x, ok := v.(bulkItem); ok {
			res = append(res, x)
		} else {
			log.Println("Unexpected item")
		}
	}
	return res
}

func (i bulkItems) snoc() (bulkItem, []bulkItem) {
	if num := len(i); num == 0 {
		return bulkItem{}, []bulkItem{}
	} else if num == 1 {
		return i[0], []bulkItem{}
	}

	return i[0], i[1:]
}

func (i bulkItems) bucketize() map[s.Channel][]bulkItem {
	res := map[s.Channel][]bulkItem{}
	for _, v := range i {
		res[v.channel] = append(res[v.channel], v)
	}
	// merge all the same channel data
	merged := map[s.Channel][]bulkItem{}
	for k, v := range res {
		if len(v) < 1 {
			continue
		}

		head, tail := bulkItems(v).snoc()
		for _, y := range tail {
			// TODO: We should validate to make sure that all the bulk items
			// have the same data, before merging
			head.members = append(head.members, y.members...)
		}
		merged[k] = []bulkItem{head}
	}
	return merged
}

func selectClusters(clusters []r.Cluster) []r.Cluster {
	// This could be a tactic (roundrobin, random, heruistic)
	index := rand.Intn(len(clusters))
	return clusters[index : index+1]
}

func scatterWrites(
	tactic Tactic,
	clusters []r.Cluster,
	fn func(r.Cluster) <-chan t.Element,
	wg *sync.WaitGroup,
	dst chan t.Element,
) error {
	return tactic(clusters, func(c r.Cluster) {
		defer wg.Done()
		for e := range fn(c) {
			dst <- e
		}
	})
}

func before(f *Farm, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr := f.instrumentation
		instr.PublishCall()
		instr.PublishKeys(1)
		instr.PublishSendTo(numSends)
	}()
	return began
}

func after(f *Farm, began time.Time) {
	go func() {
		instr := f.instrumentation
		instr.PublishDuration(time.Since(began))
	}()
}

func results(f *Farm, retrieved, returned int) {
	go func() {
		instr := f.instrumentation
		instr.PublishRetrieved(retrieved)
		instr.PublishReturned(returned)
	}()
}
