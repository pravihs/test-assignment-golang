package main

import (
	"testing"
)

func TestCache(t *testing.T) {
	cache := NewCache(3, LeastRecentlyUsed)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Put("key3", "value3")

	// Test existing key
	value, found := cache.Get("key1")
	if !found {
		t.Error("Expected key1 to be found in the cache")
	}
	if value != "value1" {
		t.Errorf("Expected value for key1 to be 'value1', got %s", value)
	}

	// Test non-existing key
	_, found = cache.Get("key4")
	if found {
		t.Error("Expected key4 not to be found in the cache")
	}

	// Test eviction
	cache.Put("key4", "value4")
	cache.Put("key5", "value5")

	// Keys key2 and key3 should have been evicted
	_, found = cache.Get("key2")
	if found {
		t.Error("Expected key2 to be evicted from the cache")
	}

	_, found = cache.Get("key3")
	if found {
		t.Error("Expected key3 to be evicted from the cache")
	}

	// Existing key should still be present
	value, found = cache.Get("key1")
	if !found {
		t.Error("Expected key1 to be found in the cache after eviction")
	}
	if value != "value1" {
		t.Errorf("Expected value for key1 to be 'value1', got %s", value)
	}
}
