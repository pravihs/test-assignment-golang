package main

import (
	"container/list"
	"sync"
	"time"
)

type EvictionStrategy int

const (
	LeastRecentlyUsed EvictionStrategy = iota
)

// CacheEntry represents an entry in the cache
type CacheEntry struct {
	key       interface{}
	value     interface{}
	timestamp time.Time
}

// Cache represents the cache module
type Cache struct {
	capacity         int
	evictionStrategy EvictionStrategy
	items            map[interface{}]*list.Element
	evictionList     *list.List
	mutex            sync.Mutex
}

func NewCache(capacity int, evictionStrategy EvictionStrategy) *Cache {
	return &Cache{
		capacity:         capacity,
		evictionStrategy: evictionStrategy,
		items:            make(map[interface{}]*list.Element),
		evictionList:     list.New(),
	}
}

// Put key-value pair into the cache
func (c *Cache) Put(key, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.evictionList.MoveToFront(element)
		element.Value.(*CacheEntry).value = value
		return
	}

	if len(c.items) >= c.capacity {
		c.evict()
	}

	element := c.evictionList.PushFront(&CacheEntry{key: key, value: value, timestamp: time.Now()})
	c.items[key] = element
}

// Get the value associated with the given key from the cache
func (c *Cache) Get(key interface{}) (value interface{}, found bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.evictionList.MoveToFront(element)
		return element.Value.(*CacheEntry).value, true
	}

	return nil, false
}

func (c *Cache) evict() {
	switch c.evictionStrategy {
	case LeastRecentlyUsed:
		c.evictLRU()
	}
}

func (c *Cache) evictLRU() {
	if element := c.evictionList.Back(); element != nil {
		c.removeElement(element)
	}
}

func (c *Cache) removeElement(element *list.Element) {
	cacheEntry := element.Value.(*CacheEntry)
	delete(c.items, cacheEntry.key)
	c.evictionList.Remove(element)
}

func main() {
	cache := NewCache(10, LeastRecentlyUsed)
	cache.Put("key1", "value1")
	value, found := cache.Get("key1")
	if found {
		println("Value for key1:", value.(string))
	} else {
		println("Key1 not found in cache")
	}
}
