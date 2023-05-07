package cache_bench

import (
	"github.com/egregors/kesh"
)

type Kesh[K comparable, V any] struct {
	c *kesh.LRUCache[K, V]
}

func NewKesh[K comparable, V any](
	size int,
) *Kesh[K, V] {
	c := kesh.NewLRUCache[K, V](size)
	return &Kesh[K, V]{
		c: c,
	}
}

func (k *Kesh[K, V]) Set(key K, value V) {
	k.c.Put(key, value)
}

func (k *Kesh[K, V]) Get(key K) (V, error) {
	return k.c.Get(key)
}

func (i *Kesh[K, V]) Del(key K) error {
	// Not implemented
	return nil
}
