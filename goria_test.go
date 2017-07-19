package goriatest

import (
	"testing"

	"github.com/oscerd/goria/gorialru"
	"github.com/oscerd/goria/goriamru"
)

func TestGoria(t *testing.T) {

	lru, err := gorialru.New("sample", 5, nil, true)
	mru, err := goriamru.New("sample", 5, nil, true)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 10; i++ {
		lru.Put(i, i)
		mru.Put(i, i)
	}

	if lru.GetName() != "sample" {
		t.Fatalf("Wrong name %v", lru.GetName())
	}

	if lru.Len() != 5 {
		t.Fatalf("Wrong len %v", lru.Len())
	}

	if mru.GetName() != "sample" {
		t.Fatalf("Wrong name %v", mru.GetName())
	}

	if mru.Len() != 5 {
		t.Fatalf("Wrong len %v", mru.Len())
	}
}
