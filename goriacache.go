package goriacache

import (
	"sync"

	"github.com/oscerd/goria/gorialru"
)

// GoriaCache is a thread-safe fixed size LRU cache.
type GoriaCache struct {
	cache *gorialru.GoriaLRU
	lock  sync.RWMutex
}

// New creates a new named GoriaLRU of the given size with stats enabled or not
func New(name string, size int, statsEnabled bool) (*GoriaCache, error) {
	return NewWithEvict(name, size, nil, statsEnabled)
}

// NewWithEvict constructs a fixed size, named cache with the given eviction and statsEnabled or not
func NewWithEvict(name string, size int, onEvicted func(key interface{}, value interface{}), statsEnabled bool) (*GoriaCache, error) {
	goriaLru, err := gorialru.NewGoriaLRU(name, size, onEvicted, statsEnabled)
	if err != nil {
		return nil, err
	}
	c := &GoriaCache{
		cache: goriaLru,
	}
	return c, nil
}

// Get looks up a key's value from the cache.
func (c *GoriaCache) Get(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.Get(key)
}

// Put add an entry in the cache
func (c *GoriaCache) Put(key interface{}, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.Put(key, value)
}

// Remove an entry with a specific value and key
func (c *GoriaCache) Remove(key interface{}, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.Remove(key, value)
}

// Replace an entry with a specific value and key with a new Value
func (c *GoriaCache) Replace(key interface{}, oldValue interface{}, newValue interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.Replace(key, oldValue, newValue)
}

// Stats get the Stats from your GoriaLRU cache
func (c *GoriaCache) Stats() gorialru.CacheStats {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cache.GetStats()
}
