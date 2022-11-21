package lorenz

import (
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/blist/selectors"
)

const (
	defaultEventCacheGCEnabled = true

	defaultEventCacheDuration   = time.Second * 30
	defaultEventCacheGCDuration = time.Minute * 5
)

type EventCache interface {
	Get(selectors.Key) (selectors.Event, bool)
	Set(selectors.Key, selectors.Event, time.Duration)
}

type noopCache struct{}

func (c *noopCache) Get(selectors.Key) (selectors.Event, bool) {
	return selectors.Event{}, false
}

func (c *noopCache) Set(selectors.Key, selectors.Event, time.Duration) {}

type eventCache struct {
	mutex *sync.Mutex
	cache map[selectors.Key]eventCacheItem
}

type eventCacheItem struct {
	event selectors.Event
	time  time.Time
}

func newEventCache() *eventCache {
	cache := &eventCache{
		mutex: &sync.Mutex{},
		cache: map[selectors.Key]eventCacheItem{},
	}
	if defaultEventCacheGCEnabled {
		go cache.gc()
	}
	return cache
}

func (c *eventCache) Get(key selectors.Key) (selectors.Event, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, ok := c.cache[key]; ok {
		if !item.time.Before(time.Now()) {
			return item.event, true
		}
		delete(c.cache, key)
	}
	return selectors.Event{}, false
}

func (c *eventCache) Set(key selectors.Key, event selectors.Event, duration time.Duration) {
	c.mutex.Lock()
	c.cache[key] = eventCacheItem{
		event: event,
		time:  time.Now().Add(duration),
	}
	c.mutex.Unlock()
}

func (c *eventCache) gc() {
	ticker := time.NewTicker(defaultEventCacheGCDuration)
	for range ticker.C {
		func(m map[selectors.Key]eventCacheItem) {
			c.mutex.Lock()

			now := time.Now()
			for k, v := range m {
				if !v.time.Before(now) {
					delete(m, k)
				}
			}

			c.mutex.Unlock()
		}(c.cache)
	}
}
