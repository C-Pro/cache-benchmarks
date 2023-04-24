package cache_bench

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/c-pro/geche"
)

const keyCardinality = 1000000

func benchmarkFuzzParallel(c geche.Geche[string, string], pb *testing.PB) {
	for pb.Next() {
		key := strconv.Itoa(rand.Intn(keyCardinality))
		r := rand.Float64()
		switch {
		case r < 0.9:
			_, _ = c.Get(key)
		case r >= 0.9 && r < 0.95:
			_ = c.Del(key)
		case r >= 0.95:
			c.Set(key, "value")
		}
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
			geche.NewRingBuffer[string, string](10000),
		},
		{
			"github.com/Code-Hex/go-generics-cache",
			NewGogLRU[string, string](ctx, time.Second, time.Second),
		},
		{
			"github.com/Yiling-J/theine-go",
			NewTheine[string, string](10000, time.Second),
		},
		{
			"github.com/jellydator/ttlcache",
			NewTTLCache[string, string](ctx, 10000, time.Second),
		},
		{
			"github.com/erni27/imcache",
			NewIMCache[string, string](time.Second),
		},
	}
	for _, c := range tab {
		b.Run(c.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				benchmarkFuzzParallel(c.imp, pb)
			})
		})
	}
}
