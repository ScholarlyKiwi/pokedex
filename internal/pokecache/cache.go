package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu       sync.Mutex
	cache    map[string]cacheEntry
	interval time.Duration
	reap     *time.Ticker
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) (*Cache, error) {
	var cache Cache
	cache.cache = make(map[string]cacheEntry)
	cache.interval = interval
	cache.reap = time.NewTicker(interval)
	go cache.reapLoop(cache.reap)

	return &cache, nil
}

func (cache *Cache) Add(key string, val []byte) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	cache.cache[key] = entry
	return nil
}

func (cache *Cache) Get(key string) (val []byte, ok bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	entry, exists := cache.cache[key]
	return entry.val, exists
}

func (cache *Cache) reapLoop(reap *time.Ticker) {
	for {
		tick := <-reap.C
		cache.mu.Lock()
		for key, value := range cache.cache {
			if tick.Compare(value.createdAt.Add(cache.interval)) >= 0 {
				delete(cache.cache, key)
			}
		}
		cache.mu.Unlock()
	}
}
