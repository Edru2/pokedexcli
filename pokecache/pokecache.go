package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mux      sync.RWMutex
	cacheMap map[string]cacheEntry
	interval time.Duration
}
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cacheMap := make(map[string]cacheEntry)
	myCache := Cache{cacheMap: cacheMap, interval: interval}
	go myCache.reapLoop()
	return &myCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	c.cacheMap[key] = cacheEntry{createdAt: time.Now(), val: val}
	c.mux.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.RLock()
	entry, ok := c.cacheMap[key]
	c.mux.RUnlock()
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for {
		limit := time.Now().Add(-c.interval)
		for key, entry := range c.cacheMap {
			if entry.createdAt.After(limit) {
				delete(c.cacheMap, key)
			}
		}
		<-ticker.C
	}
}
