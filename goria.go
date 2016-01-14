package goria

import (
	"container/list"
	"errors"
)

type EvictionCallback func(key interface{}, value interface{})

type Goria struct {
	size         int
	items        map[interface{}]*list.Element
	evictionList *list.List
	onEvict      EvictionCallback
}

type entry struct {
	key   interface{}
	value interface{}
}

func newGoria(size int, evictionC EvictionCallback) (*Goria, error) {
	if size <= 0 {
		return nil, errors.New("The Goria Cache need a positive value as size")
	}
	c := &Goria{
		size:         size,
		evictionList: list.New(),
		items:        make(map[interface{}]*list.Element),
		onEvict:      evictionC,
	}
	return c, nil
}

func (c *Goria) Put(key, value interface{}) bool {
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

func (c *Goria) Get(key interface{}) (value interface{}, exists bool) {
	if item, exists := c.items[key]; exists {
		c.evictionList.MoveToFront(item)
		return item.Value.(*entry).value, true
	}
	return
}

func (c *Goria) Remove(key interface{}) bool {
	if element, exists := c.items[key]; exists {
		c.removeElement(element)
		return true
	}
	return false
}

func (c *Goria) Keys() []interface{} {
	keys := make([]interface{}, len(c.items))
	i := 0
	for ent := c.evictionList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

func (c *Goria) removeFromTail() {
	element := c.evictionList.Back()

	if element != nil {
		c.removeElement(element)
	}
}

func (c *Goria) Len() int {
	return c.evictionList.Len()
}

func (c *Goria) removeElement(el *list.Element) {
	c.evictionList.Remove(el)
	entry := el.Value.(*entry)
	delete(c.items, entry.key)

	if c.onEvict != nil {
		c.onEvict(entry.key, entry.value)
	}
}
