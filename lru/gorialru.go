package goria

import (
	"container/list"
	"errors"
)

type EvictionCallback func(key interface{}, value interface{})

type GoriaLRU struct {
	size         int
	items        map[interface{}]*list.Element
	evictionList *list.List
	onEvict      EvictionCallback
}

type entry struct {
	key   interface{}
	value interface{}
}

func newGoriaLRU(size int, evictionC EvictionCallback) (*GoriaLRU, error) {
	if size <= 0 {
		return nil, errors.New("The Goria Cache need a positive value as size")
	}
	c := &GoriaLRU{
		size:         size,
		evictionList: list.New(),
		items:        make(map[interface{}]*list.Element),
		onEvict:      evictionC,
	}
	return c, nil
}

func (c *GoriaLRU) Put(key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		c.evictionList.MoveToFront(item)
		item.Value.(*entry).value = value
		return false
	}

	item := &entry{key, value}
	element := c.evictionList.PushFront(item)
	c.items[key] = element

	if c.evictionList.Len() > c.size {
		c.removeFromTail()
	}
	return true
}

func (c *GoriaLRU) PutIfAbsent(key, value interface{}) bool {
	var element, exists = c.items[key]
	if !exists && element == nil {
		item := &entry{key, value}
		element := c.evictionList.PushFront(item)
		c.items[key] = element

		if c.evictionList.Len() > c.size {
			c.removeFromTail()
		}
		return true
	}
	return false
}

func (c *GoriaLRU) Get(key interface{}) (value interface{}, exists bool) {
	if item, exists := c.items[key]; exists {
		c.evictionList.MoveToFront(item)
		return item.Value.(*entry).value, true
	}
	return
}

func (c *GoriaLRU) Replace(key, oldValue interface{}, newValue interface{}) bool {
	var element, exists = c.items[key]
	if exists && element != nil {
		c.evictionList.MoveToFront(element)
		element.Value.(*entry).value = newValue
		return true
	}
	return false
}

func (c *GoriaLRU) Remove(key interface{}) bool {
	if element, exists := c.items[key]; exists {
		c.removeElement(element)
		return true
	}
	return false
}

func (c *GoriaLRU) Keys() []interface{} {
	keys := make([]interface{}, len(c.items))
	i := 0
	for ent := c.evictionList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

func (c *GoriaLRU) Len() int {
	return c.evictionList.Len()
}

func (c *GoriaLRU) removeFromTail() {
	element := c.evictionList.Back()

	if element != nil {
		c.removeElement(element)
	}
}

func (c *GoriaLRU) removeElement(el *list.Element) {
	c.evictionList.Remove(el)
	entry := el.Value.(*entry)
	delete(c.items, entry.key)

	if c.onEvict != nil {
		c.onEvict(entry.key, entry.value)
	}
}
