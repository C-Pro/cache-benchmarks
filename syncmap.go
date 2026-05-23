package cache_bench

import (
	"sync"

	"github.com/c-pro/geche"
)

type SyncMap[K comparable, V any] struct {
	c    sync.Map
	zero V
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{zero: *new(V)}
}

func (s *SyncMap[K, V]) Set(key K, value V) {
	s.c.Store(key, value)
}

func (s *SyncMap[K, V]) Get(key K) (V, error) {
	v, ok := s.c.Load(key)
	if !ok {
		return s.zero, geche.ErrNotFound
	}

	return v.(V), nil
}

func (s *SyncMap[K, V]) Del(key K) error {
	s.c.Delete(key)
	return nil
}

func (s *SyncMap[K, V]) Len() int {
	// Not implemented
	return 0
}

func (s *SyncMap[K, V]) Snapshot() map[K]V {
	// Not implemented
	return nil
}

func (s *SyncMap[K, V]) SetIfPresent(K, V) (V, bool) {
	// not used in benchmark
	return s.zero, false
}

func (s *SyncMap[K, V]) Clear() {
	s.c.Clear()
}

func (s *SyncMap[K, V]) SetIfAbsent(K, V) (V, bool) {
	// not used in benchmark
	return s.zero, false
}


