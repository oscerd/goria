/*
Package Goria provides the functionality of an LRU Cache with an eye to JSR 107
*/
package goriamru

import (
	"container/list"
	"errors"
)

type EvictionCallback func(key interface{}, value interface{})

type GoriaMRU struct {
	Name         string
	Size         int
	items        map[interface{}]*list.Element
	evictionList *list.List
	onEvict      EvictionCallback
	statsEnabled bool
	stats        CacheStats
}

type CacheStats struct {
	Items     int64
	Gets      int64
	Hits      int64
	Evictions int64
	Miss      int64
}

type entry struct {
	key   interface{}
	value interface{}
}

func New(name string, size int, evictionC EvictionCallback, statsEnabled bool) (*GoriaMRU, error) {
	if size <= 0 {
		return nil, errors.New("The Goria Cache need a positive value as size")
	}
	c := &GoriaMRU{
		Name:         name,
		Size:         size,
		evictionList: list.New(),
		items:        make(map[interface{}]*list.Element),
		onEvict:      evictionC,
		statsEnabled: statsEnabled,
		stats: CacheStats{
			Items:     0,
			Evictions: 0,
			Gets:      0,
			Hits:      0,
			Miss:      0,
		},
	}
	return c, nil
}

func (c *GoriaMRU) Put(key, value interface{}) {
	if item, ok := c.items[key]; ok {
		c.evictionList.MoveToFront(item)
		item.Value.(*entry).value = value
		return
	}

	item := &entry{key, value}
	element := c.evictionList.PushFront(item)
	c.items[key] = element

	if c.evictionList.Len() > c.Size {
		c.removeFromHead()
	}

	if c.IsStatsEnabled() {
		c.stats.Items++
	}
}

func (c *GoriaMRU) PutAll(m map[interface{}]interface{}) {
	for key, value := range m {
		c.Put(key, value)
	}
}

func (c *GoriaMRU) PutIfAbsent(key, value interface{}) bool {
	var element, exists = c.items[key]
	if !exists && element == nil {

		item := &entry{key, value}
		element := c.evictionList.PushFront(item)
		c.items[key] = element

		if c.evictionList.Len() > c.Size {
			c.removeFromHead()
		}
		if c.IsStatsEnabled() {
			c.stats.Items++
		}
		return true
	}
	return false
}

func (c *GoriaMRU) Get(key interface{}) (value interface{}, exists bool) {

	if c.IsStatsEnabled() {
		c.stats.Gets++
	}
	if item, exists := c.items[key]; exists {
		c.evictionList.MoveToFront(item)
		if c.IsStatsEnabled() {
			c.stats.Hits++
		}

		return item.Value.(*entry).value, true
	}

	if c.IsStatsEnabled() {
		c.stats.Miss++
	}

	return
}

func (c *GoriaMRU) GetAll(m map[interface{}]interface{}) map[interface{}]interface{} {
	returnedMap := make(map[interface{}]interface{})

	for k := range m {
		value, exists := c.Get(k)
		if exists {
			returnedMap[k] = value
		}
	}
	return returnedMap
}

func (c *GoriaMRU) Replace(key, oldValue interface{}, newValue interface{}) bool {
	var element, exists = c.items[key]
	if exists && element.Value.(*entry).value == oldValue {
		c.evictionList.MoveToFront(element)
		element.Value.(*entry).value = newValue
		return true
	}
	return false
}

func (c *GoriaMRU) ReplaceWithKeyOnly(key, newValue interface{}) bool {
	var element, exists = c.items[key]
	if exists && element != nil {
		c.evictionList.MoveToFront(element)
		element.Value.(*entry).value = newValue
		return true
	}
	return false
}

func (c *GoriaMRU) GetAndReplace(key interface{}, newValue interface{}) interface{} {
	v, ok := c.Get(key)
	if ok {
		c.ReplaceWithKeyOnly(key, newValue)
		return v
	}
	return nil
}

func (c *GoriaMRU) RemoveWithKeyOnly(key interface{}) bool {
	if element, exists := c.items[key]; exists {
		c.removeElement(element)
		return true
	}
	return false
}

func (c *GoriaMRU) Remove(key interface{}, oldValue interface{}) bool {
	var element, exists = c.items[key]
	if exists && element.Value.(*entry).value == oldValue {
		c.removeElement(element)
		return true
	}
	return false
}

func (c *GoriaMRU) RemoveAll(m map[interface{}]interface{}) {
	for key, value := range m {
		c.Remove(key, value)
	}
}

func (c *GoriaMRU) RemoveAllWithoutParameters() {
	var keys = c.Keys()
	for i := 0; i < len(keys); i++ {
		c.RemoveWithKeyOnly(keys[i])
	}
}

func (c *GoriaMRU) GetAndRemove(key interface{}) interface{} {
	v, ok := c.Get(key)
	if ok {
		c.RemoveWithKeyOnly(key)
		return v
	}
	return nil
}

func (c *GoriaMRU) Keys() []interface{} {
	keys := make([]interface{}, len(c.items))
	i := 0
	for ent := c.evictionList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

func (c *GoriaMRU) ContainsKey(key interface{}) bool {
	for ent := c.evictionList.Back(); ent != nil; ent = ent.Prev() {
		if ent.Value.(*entry).key == key {
			return true
		}
	}
	return false
}

func (c *GoriaMRU) Len() int {
	return c.evictionList.Len()
}

func (c *GoriaMRU) GetName() string {
	return c.Name
}

func (c *GoriaMRU) IsStatsEnabled() bool {
	return c.statsEnabled
}

func (c *GoriaMRU) GetStats() CacheStats {
	return c.stats
}

func (c *GoriaMRU) removeFromHead() {
	element := c.evictionList.Front()

	if element != nil {
		c.removeElement(element)
	}
}

func (c *GoriaMRU) removeElement(el *list.Element) {
	c.evictionList.Remove(el)
	entry := el.Value.(*entry)
	delete(c.items, entry.key)

	if c.onEvict != nil {
		c.onEvict(entry.key, entry.value)
	}

	if c.IsStatsEnabled() {
		c.stats.Evictions++
		c.stats.Items--
	}
}
