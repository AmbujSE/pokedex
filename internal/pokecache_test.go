package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)

	cases := []struct {
		key string
		val []byte
	}{
		{"https://example.com", []byte("testdata")},
		{"https://example.com/path", []byte("moretestdata")},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key: %s", c.key)
			}
			if string(val) != string(c.val) {
				t.Errorf("expected value %s, got %s", string(c.val), string(val))
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)

	cache.Add("https://example.com", []byte("testdata"))
	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key before reaping")
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected key to be reaped but it was found")
	}
}
