package cache_bench

import (
	"github.com/c-pro/geche"
	lru "github.com/hashicorp/golang-lru/v2"
)

type GLRU[K comparable, V any] struct {
	c    *lru.Cache[K, V]
	zero V
}

func NewGLRU[K comparable, V any](size int) *GLRU[K, V] {
	client, err := lru.New[K, V](size)
	if err != nil {
		panic(err)
	}

	return &GLRU[K, V]{client, zero[V]()}
}

func (g *GLRU[K, V]) Get(key K) (V, error) {
	v, ok := g.c.Get(key)
	if ok {
		return v, nil
	}

	return g.zero, geche.ErrNotFound
}

func (g *GLRU[K, V]) Set(key K, value V) {
	g.c.Add(key, value)
}

func (g *GLRU[K, V]) Del(key K) error {
	g.c.Remove(key)

	return nil
}
