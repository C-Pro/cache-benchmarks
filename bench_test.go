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

type testCase struct {
	key string
	op  int
}

const (
	OPGet = iota
	OPSet
	OPDel
)

func genTestData(N int) []testCase {
	d := make([]testCase, N)
	for i := range d {
		d[i].key = strconv.Itoa(rand.Intn(keyCardinality))
		r := rand.Float64()
		switch {
		case r < 0.9:
			d[i].op = OPGet
		case r >= 0.9 && r < 0.95:
			d[i].op = OPSet
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
			geche.NewRingBuffer[string, string](100000),
		},
		{
			"github.com/Code-Hex/go-generics-cache",
			NewGogLRU[string, string](ctx, time.Second, time.Second),
		},
		{
			"github.com/Yiling-J/theine-go",
			NewTheine[string, string](100000, time.Second),
		},
		{
			"github.com/jellydator/ttlcache",
			NewTTLCache[string, string](ctx, 100000, time.Second),
		},
		{
			"github.com/erni27/imcache",
			NewIMCache[string, string](time.Second),
		},
		{
			"github.com/dgraph-io/ristretto",
			NewRistretto[string, string](100000, time.Second),
		},
		{
			"github.com/hashicorp/golang-lru/v2",
			NewGLRU[string, string](100000),
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
