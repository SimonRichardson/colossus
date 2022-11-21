package semaphore

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/typex"
)

var (
	NotFound = typex.Errorf(errors.Source, errors.NoCaseFound, "Not Found")
)

type BulkItem interface {
	Key() string
	Len() int
	Equals(BulkItem) bool
}

type queue struct {
	mutex *sync.Mutex
	items []BulkItem
}

func newQueue() *queue {
	return &queue{mutex: &sync.Mutex{}, items: []BulkItem{}}
}

func (q *queue) Add(item BulkItem) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = append(q.items, item)

	return nil
}

func (q *queue) Remove(item BulkItem) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for k, v := range q.items {
		if v.Equals(item) {
			q.items = append(q.items[:k], q.items[k+1:]...)
			return nil
		}
	}

	return NotFound
}

func (q *queue) Peek(size int) []BulkItem {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if num := len(q.items); num < size {
		size = num
	}

	var items []BulkItem
	items, q.items = q.items[:size], q.items[size:]
	return items
}

func (q *queue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	res := 0
	for _, v := range q.items {
		res += v.Len()
	}
	return res
}

type Bulk struct {
	queue   *queue
	channel chan<- []BulkItem
	quit    chan<- struct{}
	size    int
}

func NewBulk(fn func([]BulkItem), size int, timeout time.Duration) *Bulk {
	var (
		channel = make(chan []BulkItem)
		quit    = make(chan struct{})
		timer   = time.NewTicker(timeout)
		bulk    = &Bulk{
			queue:   newQueue(),
			channel: channel,
			quit:    quit,
			size:    size,
		}
	)

	go func() {
		for {
			select {
			case items := <-channel:
				go fn(items)
			case <-timer.C:
				go fn(bulk.Peek())
			case <-quit:
				timer.Stop()
				return
			}
		}
	}()

	return bulk
}

func (b *Bulk) Add(item BulkItem) error {
	if err := b.queue.Add(item); err != nil {
		return err
	}

	if num := b.queue.Len(); num > b.size {
		b.channel <- b.queue.Peek(num)
	}

	return nil
}

func (b *Bulk) Remove(item BulkItem) error {
	return b.queue.Remove(item)
}

func (b *Bulk) Peek() []BulkItem {
	return b.queue.Peek(b.size)
}

func (b *Bulk) Stop() {
	b.quit <- struct{}{}
}
