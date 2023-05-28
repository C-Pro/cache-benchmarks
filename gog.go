package cache_bench

import (
	"context"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/c-pro/geche"
)

type gogLRU[K comparable, V any] struct {
	c   *cache.Cache[K, V]
	ttl time.Duration
}

func NewGogLRU[K comparable, V any](
	ctx context.Context,
	ttl, cleanInterval time.Duration,
) *gogLRU[K, V] {
	c := cache.NewContext(ctx, cache.WithJanitorInterval[K, V](cleanInterval))
	return &gogLRU[K, V]{
		c:   c,
		ttl: ttl,
	}
}

func (g *gogLRU[K, V]) Set(key K, value V) {
	g.c.Set(key, value, cache.WithExpiration(g.ttl))
}

func (g *gogLRU[K, V]) Get(key K) (V, error) {
	val, ok := g.c.Get(key)
	if !ok {
		return val, geche.ErrNotFound
	}

	return val, nil
}

func (g *gogLRU[K, V]) Del(key K) error {
	g.c.Delete(key)

	return nil
}

func (g *gogLRU[K, V]) Len() int {
	return len(g.c.Keys())
}

func (g *gogLRU[K, V]) Snapshot() map[K]V {
	// not used in benchmark
	return nil
}
