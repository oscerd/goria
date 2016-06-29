# Goria

# Introduction

Goria is a Golang Cache that attempts to provide a Golang implementation of JSR-107 basic Cache.

Goria is based on SimpleLru by Hashicorp https://github.com/hashicorp/golang-lru

Working with Goria is simple

```golang
l, err := newGoriaLRU(128, nil)
if err != nil {
	t.Fatalf("err: %v", err)
}
for i := 0; i < 256; i++ {
	l.Put(i, i)
}
key, value := 253, 22
l.PutIfAbsent(key, value)
l.Replace(key, value, 23)
l.ReplaceWithKeyOnly(key, 24)
l.RemoveWithKeyOnly(key)
l.RemoveAll()
```
