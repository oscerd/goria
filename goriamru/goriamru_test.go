package goriamru

import (
	"fmt"
	"testing"
)

func TestGoria(t *testing.T) {

	l, err := New("sample", 5, nil, true)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 10; i++ {
		l.Put(i, i)
	}

	if l.GetName() != "sample" {
		t.Fatalf("Wrong name %v", l.GetName())
	}

	if l.Len() != 5 {
		t.Fatalf("Wrong len %v", l.Len())
	}

	for _, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || v != k {
			t.Fatalf("wrong key: %v %v", k, v)
		}
	}

	if l.GetStats().Items != 5 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 5 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	if l.Len() != 5 {
		t.Fatalf("Wrong len %v", l.Len())
	}

	key, value := 1, 1
	otherKey, oldValue, newValue := 4, 4, 47

	var result = l.PutIfAbsent(key, value)

	if l.GetStats().Items != 5 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 5 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	if result {
		t.Fatalf("key %v should be already be associated with a value", key)
	}

	result = l.PutIfAbsent(otherKey, value)

	if l.GetStats().Items != 5 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 5 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	if result {
		t.Fatalf("key %v should be associated with a value", otherKey)
	}

	l.Replace(otherKey, oldValue, newValue)

	result = l.Replace(otherKey, newValue+1, newValue+2)

	if l.GetStats().Items != 5 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 5 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	if result {
		t.Fatalf("key %v should not be replaced with a value %v, since it's current value is %v and not %v", otherKey, newValue+2, newValue, newValue+1)
	}

	result = l.ReplaceWithKeyOnly(otherKey, newValue+1)

	if !result {
		t.Fatalf("key %v should be replaced with a value %v", otherKey, newValue+1)
	}

	if l.GetStats().Items != 5 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 5 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	v, ok := l.Get(otherKey)

	if ok && v != newValue+1 {
		t.Fatalf("key %v should have a value of %v", otherKey, newValue)
	}

	res := l.evictionList.Front()

	if res.Value.(*entry).value != newValue+1 {
		fmt.Printf("%d\n", l.evictionList.Front().Value)
		fmt.Printf("%d\n", l.evictionList.Back().Value)
		t.Fatalf("key %v should have a value of %v instead has a value of %v", otherKey, newValue, res.Value.(*entry).value)
	}

	result = l.RemoveWithKeyOnly(otherKey)

	if !result {
		t.Fatalf("key %v should be removed", otherKey)
	}

	if l.GetStats().Items != 4 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 6 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	otherKey, oldValue = 1, 1
	result = l.Remove(otherKey, oldValue)

	if !result {
		t.Fatalf("key %v should be removed with a value %v", otherKey, oldValue)
	}

	if l.GetStats().Items != 3 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 7 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	if l.GetStats().Items != 3 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 7 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	otherKey = 2900
	var getAndRemoveResult = l.GetAndRemove(otherKey)

	if getAndRemoveResult != nil {
		t.Fatalf("key %v should not be removed", otherKey)
	}

	if l.GetStats().Items != 3 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 7 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	otherKey, oldValue, newValue = 3, 3, 1200
	var getAndReplaceResult = l.GetAndReplace(otherKey, newValue)

	if getAndReplaceResult != 3 {
		t.Fatalf("key %v should be replaced with an original value of %v", otherKey, oldValue)
	}

	if l.GetStats().Items != 3 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 7 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	otherKey, oldValue, newValue = 2900, 247, 1200
	getAndReplaceResult = l.GetAndReplace(otherKey, newValue)

	if getAndReplaceResult != nil {
		t.Fatalf("key %v should not be replaced", otherKey)
	}

	if l.GetStats().Items != 3 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 7 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	l.RemoveAllWithoutParameters()

	if l.GetStats().Items != 0 {
		t.Fatalf("Wrong Items stat %v", l.GetStats().Items)
	}

	if l.GetStats().Evictions != 10 {
		t.Fatalf("Wrong Evictions stat %v", l.GetStats().Evictions)
	}

	if l.GetStats().Gets != 9 {
		t.Fatalf("Wrong Gets stat %v", l.GetStats().Gets)
	}

	if l.GetStats().Hits != 7 {
		t.Fatalf("Wrong Hits stat %v", l.GetStats().Hits)
	}

	if l.GetStats().Miss != 2 {
		t.Fatalf("Wrong Miss stat %v", l.GetStats().Miss)
	}

	if l.Len() != 0 {
		t.Fatalf("Cache should be empty")
	}

}
