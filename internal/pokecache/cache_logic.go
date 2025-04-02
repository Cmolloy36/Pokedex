package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Cache    map[string]cacheEntry
	mu       *sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type BatchInfo struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func NewCache(interval time.Duration) *Cache {

	newCache := &Cache{
		Cache:    make(map[string]cacheEntry),
		mu:       &sync.Mutex{},
		interval: interval,
	}

	go newCache.reapLoop()

	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	newCacheEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	c.mu.Lock()
	c.Cache[key] = newCacheEntry

	// fmt.Println(c.Cache[key])
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	// fmt.Printf("cache key: %s\n", key)
	c.mu.Lock()
	cacheEntry, ok := c.Cache[key]
	// fmt.Println(cacheEntry)
	c.mu.Unlock()
	if !ok {
		return []byte{}, false
	}
	// fmt.Println("cache get successful")

	return cacheEntry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.Cache {
			if time.Since(entry.createdAt) >= c.interval {
				// fmt.Printf("%s deleted\n", key)
				delete(c.Cache, key)
			}
		}
		c.mu.Unlock()
	}
}
