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

	var result = l.PutIfAbsent(253, 22)

	if result {
		t.Fatalf("key %v should be already be associated with a value", 253)
	}

	result = l.PutIfAbsent(279, 22)

	if !result {
		t.Fatalf("key %v should not be associated with a value", 279)
	}
}
