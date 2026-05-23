package cache_bench

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

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

// ShardedKVCache is a sharded version of geche.KVCache
type ShardedKVCache struct {
	shards []*geche.KVCache[string, string]
	mapper geche.Mapper[string]
}

func NewShardedKVCache(numShards int) *ShardedKVCache {
	shards := make([]*geche.KVCache[string, string], numShards)
	for i := range shards {
		shards[i] = geche.NewKVCache[string, string]()
	}
	return &ShardedKVCache{
		shards: shards,
		mapper: &geche.StringMapper{},
	}
}

func (s *ShardedKVCache) Set(key, value string) {
	shardIdx := s.mapper.Map(key, len(s.shards))
	s.shards[shardIdx].Set(key, value)
}

func (s *ShardedKVCache) Get(key string) (string, error) {
	shardIdx := s.mapper.Map(key, len(s.shards))
	return s.shards[shardIdx].Get(key)
}

func (s *ShardedKVCache) ListByPrefix(prefix string) ([]string, error) {
	var res []string
	for _, shard := range s.shards {
		vals, err := shard.ListByPrefix(prefix)
		if err != nil {
			return nil, err
		}
		res = append(res, vals...)
	}
	return res, nil
}

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

// ShardedRadixCache shards RadixCache
type ShardedRadixCache struct {
	shards []*RadixCache
	mapper geche.Mapper[string]
}

func NewShardedRadixCache(numShards int) *ShardedRadixCache {
	shards := make([]*RadixCache, numShards)
	for i := range shards {
		shards[i] = NewRadixCache()
	}
	return &ShardedRadixCache{
		shards: shards,
		mapper: &geche.StringMapper{},
	}
}

func (s *ShardedRadixCache) Set(key, value string) {
	shardIdx := s.mapper.Map(key, len(s.shards))
	s.shards[shardIdx].Set(key, value)
}

func (s *ShardedRadixCache) Get(key string) (string, error) {
	shardIdx := s.mapper.Map(key, len(s.shards))
	return s.shards[shardIdx].Get(key)
}

func (s *ShardedRadixCache) ListByPrefix(prefix string) ([]string, error) {
	var res []string
	for _, shard := range s.shards {
		vals, err := shard.ListByPrefix(prefix)
		if err != nil {
			return nil, err
		}
		res = append(res, vals...)
	}
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
			name: "KVMapCacheSharded",
			factory: func() kvCache {
				return geche.NewKV[string](geche.NewSharded[string](
					func() geche.Geche[string, string] { return geche.NewMapCache[string, string]() },
					numShards,
					&geche.StringMapper{},
				))
			},
		},
		{
			name: "KVCache",
			factory: func() kvCache {
				return geche.NewKVCache[string, string]()
			},
		},
		{
			name: "KVCacheSharded",
			factory: func() kvCache {
				return NewShardedKVCache(numShards)
			},
		},
		{
			name: "ImmutableRadix",
			factory: func() kvCache {
				return NewRadixCache()
			},
		},
		{
			name: "ImmutableRadixSharded",
			factory: func() kvCache {
				return NewShardedRadixCache(numShards)
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


