package cache_bench

import (
	"time"

	"github.com/Yiling-J/theine-go"
	"github.com/c-pro/geche"
)

type theineCache[K comparable, V any] struct {
	c   *theine.Cache[K, V]
	ttl time.Duration
}

func NewTheine[K comparable, V any](
	size int64,
	ttl time.Duration,
) *theineCache[K, V] {
	c, err := theine.NewBuilder[K, V](size).Build()
	if err != nil {
		panic(err)
	}

	return &theineCache[K, V]{
		c:   c,
		ttl: ttl,
	}
}

func (t *theineCache[K, V]) Set(key K, value V) {
	t.c.SetWithTTL(key, value, 1, t.ttl)
}

func (t *theineCache[K, V]) Get(key K) (V, error) {
	val, ok := t.c.Get(key)
	if !ok {
		return val, geche.ErrNotFound
	}

	return val, nil
}

func (t *theineCache[K, V]) Del(key K) error {
	t.c.Delete(key)

	return nil
}

func (t *theineCache[K, V]) Len() int {
	return t.c.Len()
}

func (t *theineCache[K, V]) Snapshot() map[K]V {
	// not used in benchmark
	return nil
}
