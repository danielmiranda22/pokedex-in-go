package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	create_at time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.Mutex    // mutex for concurrent access to the cache - goroutines
	interval time.Duration // how long entries live
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}

	// start reaper in background — runs for the lifetime of the program
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		create_at: time.Now(),
		val:       val,
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.entries {
			if time.Since(v.create_at) > interval {
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}
