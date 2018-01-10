package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	var found bool

	cache := NewCache(50 * time.Millisecond)
	cache.Set("a", 1, DefaultExpiration)
	cache.Set("b", 2, NoExpiration)
	cache.Set("c", 3, 20*time.Millisecond)
	cache.Set("d", 4, 30*time.Millisecond)

	<-time.After(25 * time.Millisecond) // elapsed: ~25ms

	_, found = cache.Get("c")
	if found {
		t.Error("Found c in cache when it should have been evicted.")
	}

	_, found = cache.Get("d")
	if !found {
		t.Error("Didn't find d in cache when it should be there for 5 more ms.")
	}

	<-time.After(10 * time.Millisecond) // elapsed: ~35ms

	_, found = cache.Get("d")
	if found {
		t.Error("Found d in cache when it should have been evicted.")
	}

	<-time.After(20 * time.Millisecond) // elapsed: ~55ms

	_, found = cache.Get("a")
	if found {
		t.Error("Found a in cache when it should have been evicted.")
	}

	_, found = cache.Get("b")
	if !found {
		t.Error("Didn't find b in cache even though it should never be evicted.")
	}
}
