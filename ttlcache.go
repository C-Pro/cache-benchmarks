package cache_bench

import (
	"context"
	"time"

	"github.com/c-pro/geche"
	cache "github.com/jellydator/ttlcache/v3"
)

type TTLCache[K comparable, V any] struct {
	c    *cache.Cache[K, V]
	zero V
}

func zero[V any]() V {
	var z V
	return z
}

func NewTTLCache[K comparable, V any](
	ctx context.Context,
	size uint64,
	ttl time.Duration,
) *TTLCache[K, V] {
	c := cache.New(
		cache.WithCapacity[K, V](size),
		cache.WithTTL[K, V](ttl),
	)

	go c.Start()

	go func() {
		<-ctx.Done()
		c.Stop()
	}()
	return &TTLCache[K, V]{
		c:    c,
		zero: zero[V](),
	}
}

func (t *TTLCache[K, V]) Set(key K, value V) {
	t.c.Set(key, value, cache.DefaultTTL)
}

func (t *TTLCache[K, V]) Get(key K) (V, error) {
	val := t.c.Get(key)
	if val == nil || val.IsExpired() {
		return t.zero, geche.ErrNotFound
	}

	return val.Value(), nil
}

func (t *TTLCache[K, V]) Del(key K) error {
	t.c.Delete(key)

	return nil
}
