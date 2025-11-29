package cache_bench

import (
	"context"
	"fmt"
	"maps"
	"math"
	"math/rand"
	"runtime"
	"slices"
	"testing"
	"time"

	"github.com/c-pro/geche"
	"github.com/erni27/imcache"
)

const keyCardinality = 1000000

var numShards = 1 << (int(math.Log2(float64(runtime.NumCPU()))) + 1)

type testCase struct {
	key string
	op  int
}

const (
	OPGet = iota
	OPSet
	OPDel
)

func genRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.Intn(26) + int(byte('a')))
	}
	return string(b)
}

func genTestData(N int) []testCase {
	// Generate composite keys with common real-life pattern.
	// {keyType}:{userID}:{objectID}
	// This is mostly irrelevant for map-backed caches, but makes a world of difference
	// for trie-based KV cache.
	var (
		numKeyTypes = 20
		numUsers    = N / 1000
		// Number of distinct keys controls how much hits/misses we have in the benchmark.
		// When number of distinct keys >= N, we have mostly misses.
		distinctKeys = N / 10
	)
	keyTypesMap := make(map[string]struct{}, numKeyTypes)
	for len(keyTypesMap) < numKeyTypes {
		keyTypesMap[genRandomString(5)] = struct{}{}
	}
	keyTypes := slices.Collect(maps.Keys(keyTypesMap))

	usersMap := make(map[string]struct{}, numUsers)
	for len(usersMap) < numUsers {
		usersMap[genRandomString(16)] = struct{}{}
	}
	users := slices.Collect(maps.Keys(usersMap))

	keysMap := make(map[string]struct{}, distinctKeys)
	for len(keysMap) < distinctKeys {
		key := fmt.Sprintf("%s:%s:%s",
			keyTypes[rand.Intn(len(keyTypes))],
			users[rand.Intn(len(users))],
			genRandomString(16),
		)
		keysMap[key] = struct{}{}
	}

	keys := slices.Collect(maps.Keys(keysMap))

	d := make([]testCase, N)
	for i := range d {
		d[i].key = keys[rand.Intn(len(keys))]
		r := rand.Float64()
		switch {
		// Write heavy, because with read heavy, most of the reads would be misses.
		case r < 0.7:
			d[i].op = OPSet
		case r >= 0.7 && r < 0.95:
			d[i].op = OPGet
		case r >= 0.95:
			d[i].op = OPDel
		}
	}

	return d
}

func benchmarkFuzzParallel(
	c geche.Geche[string, string],
	testData []testCase,
	pb *testing.PB,
) {
	i := 0
	for pb.Next() {
		switch testData[i].op {
		case OPGet:
			_, _ = c.Get(testData[i].key)
		case OPSet:
			c.Set(testData[i].key, "value")
		case OPDel:
			_ = c.Del(testData[i].key)
		}
		i = (i + 1) % len(testData)
	}
}

// BenchmarkEverything performs different operations randomly.
// Ratio for get/set/del is 90/5/5
func BenchmarkEverythingParallel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tab := []struct {
		name string
		imp  geche.Geche[string, string]
	}{
		{
			"MapCache",
			geche.NewMapCache[string, string](),
		},
		{
			"MapTTLCache",
			geche.NewMapTTLCache[string, string](ctx, time.Second, time.Second),
		},
		{
			"RingBuffer",
			geche.NewRingBuffer[string, string](1000000),
		},
		{
			"ShardedMapCache",
			geche.NewSharded[string](
				func() geche.Geche[string, string] { return geche.NewMapCache[string, string]() },
				numShards,
				&geche.StringMapper{},
			),
		},
		{
			"ShardedMapTTLCache",
			geche.NewSharded[string](
				func() geche.Geche[string, string] {
					return geche.NewMapTTLCache[string, string](ctx, time.Second, time.Second)
				},
				numShards,
				&geche.StringMapper{},
			),
		},
		{
			"ShardedRingBuffer",
			geche.NewSharded[string](
				func() geche.Geche[string, string] { return geche.NewRingBuffer[string, string](1000000/numShards + 1) },
				numShards,
				&geche.StringMapper{},
			),
		},
		{
			"github.com/Code-Hex/go-generics-cache",
			NewGogLRU[string, string](ctx, time.Second, time.Second),
		},
		{
			"github.com/Yiling-J/theine-go",
			NewTheine[string, string](1000000, time.Second),
		},
		{
			"github.com/jellydator/ttlcache",
			NewTTLCache[string, string](ctx, 1000000, time.Second),
		},
		{
			"github.com/erni27/imcache",
			NewIMCache[string, string](time.Second, imcache.DefaultStringHasher64{}, numShards),
		},
		{
			"github.com/dgraph-io/ristretto",
			NewRistretto[string, string](1000000, time.Second),
		},
		{
			"github.com/hashicorp/golang-lru/v2",
			NewGLRU[string, string](1000000),
		},
		{
			"github.com/egregors/kesh",
			NewKesh[string, string](1000000),
		},
		{
			"KVMapCache",
			geche.NewKV[string](geche.NewMapCache[string, string]()),
		},
		{
			"KVCache",
			geche.NewKVCache[string, string](),
		},
		{
			"ShardedKVCache",
			geche.NewSharded(
				func() geche.Geche[string, string] { return geche.NewKVCache[string, string]() },
				numShards,
				&geche.StringMapper{},
			),
		},
	}
	data := genTestData(10_000_000)
	b.ResetTimer()
	for _, c := range tab {
		b.Run(c.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				benchmarkFuzzParallel(c.imp, data, pb)
			})
		})
	}
}
