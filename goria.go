package goria

import (
	"errors"
)

type EvictionCallback func(key interface{}, value interface{})

type Goria struct {
	size    int
	items   map[interface{}]interface{}
	onEvict EvictionCallback
}

func newGoria(size int, evictionC EvictionCallback) (*Goria, error) {
	if size <= 0 {
		return nil, errors.New("The Goria Cache need a positive value as size")
	}
	c := &Goria{
		size:    size,
		items:   make(map[interface{}]interface{}),
		onEvict: evictionC,
	}
	return c, nil
}

func (c *Goria) Add(key, value interface{}) bool {
	c.items[key] = value
	return true
}

func (c *Goria) Len() int {
	return len(c.items)
}
