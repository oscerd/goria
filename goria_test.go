package goria

import (
	"fmt"
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
		l.Add(i, i)
	}

	if l.Len() != 256 {
		t.Fatalf("bad len: %v", l.Len())
	}

	fmt.Printf("%v\n", l.Len())
}
