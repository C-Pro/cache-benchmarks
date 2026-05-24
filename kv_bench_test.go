package cache_bench

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/armon/go-radix"
	"github.com/c-pro/geche"
	"github.com/google/uuid"
	iradix "github.com/hashicorp/go-immutable-radix"
)

var (
	testKeys     []string
	testPrefixes []string
)

func init() {
	// Generate 100,000 keys
	testKeys = make([]string, 100000)
	for i := 0; i < len(testKeys); i++ {
		category := fmt.Sprintf("prefix-%02d-", i%100)
		testKeys[i] = category + uuid.NewString()
	}

	// Generate some prefixes for querying
	testPrefixes = []string{
		"prefix-00-",   // matches 1,000
		"prefix-25-",   // matches 1,000
		"prefix-50-",   // matches 1,000
		"prefix-75-",   // matches 1,000
		"prefix-0",     // matches 10,000
		"prefix-5",     // matches 10,000
		"prefix-",      // matches 100,000
		"nonexistent-", // matches 0
	}
}

type kvCache interface {
	Set(key, value string)
	Get(key string) (string, error)
	ListByPrefix(prefix string) ([]string, error)
}

// ----------------- Wrappers -----------------

// RadixCache wraps hashicorp/go-immutable-radix to be mutable and thread-safe
type RadixCache struct {
	mux sync.RWMutex
	t   *iradix.Tree
}

func NewRadixCache() *RadixCache {
	return &RadixCache{
		t: iradix.New(),
	}
}

func (r *RadixCache) Set(key, value string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.t, _, _ = r.t.Insert([]byte(key), value)
}

func (r *RadixCache) Get(key string) (string, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	val, ok := r.t.Get([]byte(key))
	if !ok {
		return "", geche.ErrNotFound
	}
	return val.(string), nil
}

func (r *RadixCache) ListByPrefix(prefix string) ([]string, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	var res []string
	r.t.Root().WalkPrefix([]byte(prefix), func(k []byte, v interface{}) bool {
		res = append(res, v.(string))
		return false
	})
	return res, nil
}

// ArmonRadixCache wraps armon/go-radix to be thread-safe
type ArmonRadixCache struct {
	mux sync.RWMutex
	t   *radix.Tree
}

func NewArmonRadixCache() *ArmonRadixCache {
	return &ArmonRadixCache{
		t: radix.New(),
	}
}

func (r *ArmonRadixCache) Set(key, value string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.t.Insert(key, value)
}

func (r *ArmonRadixCache) Get(key string) (string, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	val, ok := r.t.Get(key)
	if !ok {
		return "", geche.ErrNotFound
	}
	return val.(string), nil
}

func (r *ArmonRadixCache) ListByPrefix(prefix string) ([]string, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	var res []string
	r.t.WalkPrefix(prefix, func(s string, v interface{}) bool {
		res = append(res, v.(string))
		return false
	})
	return res, nil
}

// ----------------- Benchmark Registrations -----------------

func getCaches() []struct {
	name    string
	factory func() kvCache
} {
	return []struct {
		name    string
		factory func() kvCache
	}{
		{
			name: "KVMapCache",
			factory: func() kvCache {
				return geche.NewKV[string](geche.NewMapCache[string, string]())
			},
		},
		{
			name: "KVCache",
			factory: func() kvCache {
				return geche.NewKVCache[string, string]()
			},
		},
		{
			name: "ImmutableRadix",
			factory: func() kvCache {
				return NewRadixCache()
			},
		},
		{
			name: "ArmonRadix",
			factory: func() kvCache {
				return NewArmonRadixCache()
			},
		},
	}
}

func BenchmarkKV_Insert(b *testing.B) {
	for _, tc := range getCaches() {
		b.Run(tc.name, func(b *testing.B) {
			c := tc.factory()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := rand.Intn(len(testKeys))
				for pb.Next() {
					c.Set(testKeys[i%len(testKeys)], "value")
					i++
				}
			})
		})
	}
}

func BenchmarkKV_Get(b *testing.B) {
	for _, tc := range getCaches() {
		b.Run(tc.name, func(b *testing.B) {
			c := tc.factory()
			for _, k := range testKeys {
				c.Set(k, "value")
			}
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := rand.Intn(len(testKeys))
				for pb.Next() {
					_, _ = c.Get(testKeys[i%len(testKeys)])
					i++
				}
			})
		})
	}
}

func BenchmarkKV_ListByPrefix(b *testing.B) {
	for _, tc := range getCaches() {
		b.Run(tc.name, func(b *testing.B) {
			c := tc.factory()
			for _, k := range testKeys {
				c.Set(k, "value")
			}
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := rand.Intn(len(testPrefixes))
				for pb.Next() {
					_, _ = c.ListByPrefix(testPrefixes[i%len(testPrefixes)])
					i++
				}
			})
		})
	}
}
