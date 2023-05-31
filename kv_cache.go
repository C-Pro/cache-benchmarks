package cache_bench

import (
	"sync"

	"github.com/c-pro/geche"
)

type trieCacheNode[V any] struct {
	c        byte
	next     [256]*trieCacheNode[V]
	terminal bool
	value    V
}

type KVCache[V any] struct {
	trie *trieCacheNode[V]
	mux  sync.RWMutex
	zero V
	len  int
}

// KVCache is a cache that supports prefix search on keys using ListByPrefix function.
func NewKVCache[V any]() *KVCache[V] {
	kv := KVCache[V]{
		trie: new(trieCacheNode[V]),
		zero: zero[V](),
	}

	return &kv
}

// Set key-value pair while updating the trie.
func (kv *KVCache[V]) Set(key string, value V) {
	kv.mux.Lock()
	defer kv.mux.Unlock()

	created := false
	node := kv.trie
	for i := 0; i < len(key); i++ {
		next := node.next[key[i]]
		if next == nil {
			next = &trieCacheNode[V]{
				c: key[i],
			}
			node.next[key[i]] = next
			created = true
		}

		node = next
	}

	node.terminal = true
	node.value = value
	if created {
		kv.len++
	}
}

func (kv *KVCache[V]) dfs(node *trieCacheNode[V], prefix []byte) ([]V, error) {
	res := []V{}
	if node.terminal {
		res = append(res, node.value)
	}

	for i := 0; i < len(node.next); i++ {
		if node.next[i] != nil {
			next := node.next[i]
			nextRes, err := kv.dfs(next, append(prefix, next.c))
			if err != nil {
				return nil, err
			}
			res = append(res, nextRes...)
		}
	}

	return res, nil
}

func (kv *KVCache[V]) ListByPrefix(prefix string) ([]V, error) {
	kv.mux.RLock()
	defer kv.mux.RUnlock()

	node := kv.trie
	for i := 0; i < len(prefix); i++ {
		next := node.next[prefix[i]]
		if next == nil {
			return nil, nil
		}
		node = next
	}

	return kv.dfs(node, []byte(prefix))
}

// Get value by key from the underlying cache.
func (kv *KVCache[V]) Get(key string) (V, error) {
	node := kv.trie
	for i := 0; i < len(key); i++ {
		if node.next[i] == nil {
			return kv.zero, geche.ErrNotFound
		}
		node = node.next[key[i]]
	}

	return node.value, nil
}

// Del key from the cache.
func (kv *KVCache[V]) Del(key string) error {
	kv.mux.Lock()
	defer kv.mux.Unlock()

	node := kv.trie
	var prev *trieCacheNode[V]
	for i := 0; i < len(key); i++ {
		next := node.next[key[i]]
		if next == nil {
			return nil
		}

		prev = node
		node = next
	}

	node.terminal = false

	empty := true
	for i := 0; i < len(node.next); i++ {
		if node.next[i] != nil {
			empty = false
			break
		}
	}

	if empty {
		prev.next[node.c] = nil
	}

	kv.len--

	return nil
}

// Snapshot returns a shallow copy of the cache data.
func (kv *KVCache[V]) Snapshot() map[string]V {
	// TODO
	return nil
}

// Len returns total number of elements in the underlying caches.
func (kv *KVCache[V]) Len() int {
	kv.mux.RLock()
	defer kv.mux.RUnlock()
	return kv.len
}
