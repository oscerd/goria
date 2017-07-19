# Goria

# Introduction

Goria is a Golang Cache that attempts to provide a Golang implementation of JSR-107 basic Cache.

Goria is based on SimpleLru by Hashicorp https://github.com/hashicorp/golang-lru

Working with a GoriaCache is simple

```golang
cache, err := NewWithEvict("sample", 128, nil, true)
if err != nil {
	t.Fatalf("err: %v", err)
}
for i := 0; i < 256; i++ {
	cache.Put(i, i)
}
v, ok := cache.Get(255)
keyValuesSet := map[interface{}]interface{}{
	10: 10,
	20: 20,
	30: 30,
	40: 40,
}
returnedKeyValuesSet := make(map[interface{}]interface{})
returnedKeyValuesSet = cache.GetAll(keyValuesSet)
keyValuesSet = map[interface{}]interface{}{
	10: 11,
	20: 22,
	30: 33,
	40: 44,
}
cache.PutAll(keyValuesSet)
if cache.IsStatsEnabled() {
	fmt.Printf("Evictions %v\n", cache.Stats().Evictions)
	fmt.Printf("Items %v\n", cache.Stats().Items)
}
```

Currently we have LRU and MRU policy caches.
