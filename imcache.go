package cache_bench

import (
	"time"

	"github.com/c-pro/geche"
	"github.com/erni27/imcache"
)

type IMCache[K comparable, V any] struct {
	c *imcache.Sharded[K, V]
}

func NewIMCache[K comparable, V any](
	ttl time.Duration,
	hasher imcache.Hasher64[K],
	numShards int,
) *IMCache[K, V] {
	c := imcache.NewSharded(
		numShards,
		hasher,
		imcache.WithDefaultExpirationOption[K, V](ttl),
	)

	return &IMCache[K, V]{
		c: c,
	}
}

func (i *IMCache[K, V]) Set(key K, value V) {
	i.c.Set(key, value, imcache.WithDefaultExpiration())
}

func (i *IMCache[K, V]) Get(key K) (V, error) {
	val, ok := i.c.Get(key)
	if !ok {
		return val, geche.ErrNotFound
	}

	return val, nil
}

func (i *IMCache[K, V]) Del(key K) error {
	i.c.Remove(key)

	return nil
}

func (i *IMCache[K, V]) Len() int {
	return i.c.Len()
}

func (i *IMCache[K, V]) Snapshot() map[K]V {
	// not used in benchmark
	return nil
}
