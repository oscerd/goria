package goria

import (
	"testing"
)

func TestGoria(t *testing.T) {
	onEvicted := func(k interface{}, v interface{}) {
		if k != v {
			t.Fatalf("Evict values not equal (%v!=%v)", k, v)
		}
	}

	l, err := newGoria(128, onEvicted)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		l.Put(i, i)
	}

	if l.Len() != 128 {
		t.Fatalf("Wrong len %v", l.Len())
	}

	for i, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || v != k || v != i+128 {
			t.Fatalf("wrong key: %v", k)
		}
	}

	key, value := 253, 22
	otherKey, oldValue, newValue := 279, 22, 47

	var result = l.PutIfAbsent(key, value)

	if result {
		t.Fatalf("key %v should be already be associated with a value", key)
	}

	result = l.PutIfAbsent(otherKey, value)

	if !result {
		t.Fatalf("key %v should not be associated with a value", otherKey)
	}

	l.Replace(otherKey, oldValue, newValue)

	v, ok := l.Get(otherKey)

	if ok && v != newValue {
		t.Fatalf("key %v should have a value of %v", otherKey, newValue)
	}

	res := l.evictionList.Front()

	if res.Value.(*entry).value != newValue {
		t.Fatalf("key %v should have a value of 47 instead has a value of %v", otherKey, res.Value.(*entry).value)
	}
}
