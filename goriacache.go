package goriacachelru

import (
	"sync"

	"github.com/oscerd/goria/gorialru"
)

// GoriaCacheLRU is a thread-safe fixed size LRU cache.
type GoriaCacheLRU struct {
	cache *gorialru.GoriaLRU
	lock  sync.RWMutex
}

// New creates a new named GoriaLRU of the given size with stats enabled or not
func New(name string, size int, statsEnabled bool) (*GoriaCacheLRU, error) {
	return NewWithEvict(name, size, nil, statsEnabled)
}

// NewWithEvict constructs a fixed size, named cache with the given eviction and statsEnabled or not
func NewWithEvict(name string, size int, onEvicted func(key interface{}, value interface{}), statsEnabled bool) (*GoriaCacheLRU, error) {
	goriaLru, err := gorialru.New(name, size, onEvicted, statsEnabled)
	if err != nil {
		return nil, err
	}
	c := &GoriaCacheLRU{
		cache: goriaLru,
	}
	return c, nil
}

// Get looks up a key's value from the cache.
func (c *GoriaCacheLRU) Get(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.Get(key)
}

// GetAll return a set of Key/Value pairs.
func (c *GoriaCacheLRU) GetAll(m map[interface{}]interface{}) map[interface{}]interface{} {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.GetAll(m)
}

// Put add an entry in the cache
func (c *GoriaCacheLRU) Put(key interface{}, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.Put(key, value)
}

// PutAll return a set of Key/Value pairs.
func (c *GoriaCacheLRU) PutAll(m map[interface{}]interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.PutAll(m)
}

// Remove an entry with a specific value and key
func (c *GoriaCacheLRU) Remove(key interface{}, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.Remove(key, value)
}

// Replace an entry with a specific value and key with a new Value
func (c *GoriaCacheLRU) Replace(key interface{}, oldValue interface{}, newValue interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cache.Replace(key, oldValue, newValue)
}

// Stats get the Stats from your GoriaLRU cache
func (c *GoriaCacheLRU) Stats() gorialru.CacheStats {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cache.GetStats()
}

// IsStatsEnabled is a method used to check if stats are enabled or not
func (c *GoriaCacheLRU) IsStatsEnabled() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cache.IsStatsEnabled()
}
