package pokecache

import (
	"time"
	"sync"
	"fmt"	
)

// Initialises a new cache
func NewCache(duration int) *Cache {
	newCache := Cache{}
	newCache.Cache = make(map[string]cacheEntry,0)
	go newCache.reaper(5)
	return &newCache
}

// go reaper looks after the cache deleting old data
func (c *Cache) reaper(duration_seconds int) {

	// tick tock cache!
	duration := time.Duration(duration_seconds) * time.Second
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for t := range ticker.C{
		fmt.Printf("Checked Cache at: %v", t)
		for key, value := range(c.Cache) {
			if time.Since(value.createdAt) > duration {
				// delete cache item
				delete(c.Cache, key)
			}
		}
	}
}

type Cache struct {
	Cache map[string]cacheEntry
	mu   sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	value []byte
}

// get the value from the cache
func (c *Cache) Get(key string) ([]byte, bool){
	c.mu.Lock()
    defer c.mu.Unlock()
	data, ok := c.Cache[key]
	if ok == false {
		return nil, false
	}
	return data.value, true
}

// stores data in the cache
func (c *Cache) Add(key string, data []byte) {
	c.mu.Lock()
	c.Cache[key] = cacheEntry{
		createdAt: time.Now(),
		value: data,
	}
    defer c.mu.Unlock()
}

