package cache_bench

import (
	"time"

	"github.com/c-pro/geche"
	"github.com/dgraph-io/ristretto"
)

type Ristretto[K comparable, V any] struct {
	c    *ristretto.Cache
	zero V
	ttl  time.Duration
}

func NewRistretto[K comparable, V any](size int, ttl time.Duration) *Ristretto[K, V] {
	client, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: int64(size * 10),
		MaxCost:     int64(size),
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}

	return &Ristretto[K, V]{client, zero[V](), ttl}
}

func (r *Ristretto[K, V]) Get(key K) (V, error) {
	v, ok := r.c.Get(key)
	if ok {
		return v.(V), nil
	}

	return r.zero, geche.ErrNotFound
}

func (r *Ristretto[K, V]) Set(key K, value V) {
	r.c.SetWithTTL(key, value, 1, r.ttl)
}

func (r *Ristretto[K, V]) Del(key K) error {
	r.c.Del(key)

	return nil
}

func (r *Ristretto[K, V]) Len() int {
	return int(r.c.Metrics.KeysAdded())
}

func (r *Ristretto[K, V]) Snapshot() map[K]V {
	// not used in benchmark
	return nil
}
