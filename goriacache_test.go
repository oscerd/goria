package goriacache

import "testing"

func TestGoriaCache(t *testing.T) {
	cache, err := New("sample", 256, true)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		cache.Put(i, i)
	}

	if cache.Stats().Items != 256 {
		t.Fatalf("Wrong Items stats %v", cache.Stats().Items)
	}

	v, ok := cache.Get(255)
	if v != 255 && !ok {
		t.Fatalf("255 should be in the cache with key 255")
	}
}
