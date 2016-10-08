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

	if cache.IsStatsEnabled() {
		if cache.Stats().Items != 256 {
			t.Fatalf("Wrong Items stats %v", cache.Stats().Items)
		}
	}

	v, ok := cache.Get(255)
	if v != 255 && !ok {
		t.Fatalf("255 should be in the cache with key 255")
	}

	keyValuesSet := map[interface{}]interface{}{
		10: 10,
		20: 20,
		30: 30,
		40: 40,
	}

	returnedKeyValuesSet := make(map[interface{}]interface{})
	returnedKeyValuesSet = cache.GetAll(keyValuesSet)

	for k, v := range keyValuesSet {
		if returnedKeyValuesSet[k] != v {
			t.Fatalf("key %v should have value %v", k, v)
		}
	}

	keyValuesSet = map[interface{}]interface{}{
		10: 11,
		20: 22,
		30: 33,
		40: 44,
	}

	cache.PutAll(keyValuesSet)

	returnedKeyValuesSet = cache.GetAll(keyValuesSet)

	for k, v := range keyValuesSet {
		if returnedKeyValuesSet[k] != v {
			t.Fatalf("key %v should have value %v", k, v)
		}
	}

	ok = cache.Replace(40, 44, 40)
	if !ok {
		t.Fatalf("err: %v", err)
	}

	v, ok = cache.Get(40)
	if !ok && v != 40 {
		t.Fatalf("key %v should have value %v", 40, 40)
	}

	ok = cache.Remove(40, 40)
	if !ok {
		t.Fatalf("err: %v", err)
	}
}
