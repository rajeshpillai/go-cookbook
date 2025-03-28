// Question 4: Building a Thread-Safe In-Memory Cache
// Design and implement a thread-safe in-memory cache in Go. The cache should provide basic operations (Get, Set, Delete) and support expiration of items. Describe your design choices regarding concurrency control and eviction strategy."

package main

import (
	"fmt"
	"sync"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration int64
}

type Cache struct {
	items map[string]cacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]cacheItem),
		ttl:   ttl,
	}
	go c.startEviction()
	return c
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl).UnixNano(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.items[key]
	if !exists || time.Now().UnixNano() > item.expiration {
		return nil, false
	}
	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) startEviction() {
	ticker := time.NewTicker(c.ttl)
	for {
		<-ticker.C
		c.mu.Lock()
		now := time.Now().UnixNano()
		for k, item := range c.items {
			if now > item.expiration {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

func main() {
	cache := NewCache(2 * time.Second)
	cache.Set("foo", "bar")
	if value, found := cache.Get("foo"); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not Found")
	}
	time.Sleep(3 * time.Second)
	if value, found := cache.Get("foo"); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not Found after expiration")
	}
}

// Review Points:
// Concurrency: Uses sync.RWMutex to ensure safe concurrent access.
// Expiration Strategy: Implements a background goroutine to periodically evict expired items.
// Design Trade-offs: Candidate can discuss alternative eviction strategies (e.g., LRU) and improvements like using third-party caching libraries.
