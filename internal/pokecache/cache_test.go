package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	cache := NewCache(time.Minute)
	cache.Add("key", []byte("value"))

	val, ok := cache.Get("key")
	if !ok {
		t.Fatal("expected to find key in cache")
	}
	if string(val) != "value" {
		t.Errorf("got %q, want %q", val, "value")
	}
}

func TestCacheReap(t *testing.T) {
	interval := 5 * time.Millisecond
	cache := NewCache(interval)
	cache.Add("key", []byte("value"))

	time.Sleep(interval * 2)

	_, ok := cache.Get("key")
	if ok {
		t.Fatal("expected key to be reaped from cache")
	}
}
