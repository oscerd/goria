package goria

import "testing"

func TestGoria(t *testing.T) {

	l, err := newGoriaLRU("sample", 128, nil)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		l.Put(i, i)
	}

	if l.Name() != "sample" {
		t.Fatalf("Wrong name %v", l.Name())
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

	result = l.Replace(otherKey, newValue+1, newValue+2)

	if result {
		t.Fatalf("key %v should not be replaced with a value %v, since it's current value is %v and not %v", otherKey, newValue+2, newValue, newValue+1)
	}

	result = l.ReplaceWithKeyOnly(otherKey, newValue+1)

	if !result {
		t.Fatalf("key %v should be replaced with a value %v", otherKey, newValue+1)
	}

	v, ok := l.Get(otherKey)

	if ok && v != newValue+1 {
		t.Fatalf("key %v should have a value of %v", otherKey, newValue)
	}

	res := l.evictionList.Front()

	if res.Value.(*entry).value != newValue+1 {
		t.Fatalf("key %v should have a value of %v instead has a value of %v", otherKey, newValue, res.Value.(*entry).value)
	}

	result = l.RemoveWithKeyOnly(otherKey)

	if !result {
		t.Fatalf("key %v should be removed", otherKey)
	}

	otherKey, oldValue = 252, 252
	result = l.Remove(otherKey, oldValue)

	if !result {
		t.Fatalf("key %v should be removed with a value %v", otherKey, oldValue)
	}

	otherKey = 248
	var getAndRemoveResult = l.GetAndRemove(otherKey)

	if getAndRemoveResult != 248 {
		t.Fatalf("key %v should be removed with a value %v", otherKey)
	}

	otherKey = 2900
	getAndRemoveResult = l.GetAndRemove(otherKey)

	if getAndRemoveResult != nil {
		t.Fatalf("key %v should not be removed")
	}

	otherKey, oldValue, newValue = 247, 247, 1200
	var getAndReplaceResult = l.GetAndReplace(otherKey, newValue)

	if getAndReplaceResult != 247 {
		t.Fatalf("key %v should be replaced with an original value of %v", otherKey, oldValue)
	}

	otherKey, oldValue, newValue = 2900, 247, 1200
	getAndReplaceResult = l.GetAndReplace(otherKey, newValue)

	if getAndReplaceResult != nil {
		t.Fatalf("key %v should not be replaced", otherKey)
	}

	commits := map[interface{}]interface{}{
		253: 24,
		267: 22,
		280: 21,
		281: 23,
	}
	var newKey = 253
	newValue = 24
	l.PutAll(commits)
	v, ok = l.Get(newKey)

	if v != newValue {
		t.Fatalf("key %v should have value %v", newKey, newValue)
	}

	l.RemoveAll(commits)

	v, ok = l.Get(newKey)

	if ok {
		t.Fatalf("key %v shouldn't be in the cache", newKey)
	}

	l.PutAll(commits)

	l.RemoveAllWithoutParameters()

	v, ok = l.Get(newKey)

	if ok {
		t.Fatalf("key %v shouldn't be in the cache", newKey)
	}

	if l.Len() != 0 {
		t.Fatalf("Cache should be empty")
	}
}