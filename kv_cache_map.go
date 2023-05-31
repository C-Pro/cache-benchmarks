package cache_bench

import (
	"sync"

	"github.com/c-pro/geche"
)

type trieCacheNodeMap[V any] struct {
	c     byte
	next  map[byte]*trieCacheNodeMap[V]
	value *V
}

type KVCacheMap[V any] struct {
	trie *trieCacheNodeMap[V]
	mux  sync.RWMutex
	zero V
	len  int
}

// KVCacheMap is a cache that supports prefix search on keys using ListByPrefix function.
func NewKVCacheMap[V any]() *KVCacheMap[V] {
	kv := KVCacheMap[V]{
		trie: &trieCacheNodeMap[V]{},
		zero: zero[V](),
	}

	return &kv
}

// Set key-value pair while updating the trie.
func (kv *KVCacheMap[V]) Set(key string, value V) {
	kv.mux.Lock()
	defer kv.mux.Unlock()

	node := kv.trie
	prev := node
	for i:=0; i < len(key); i++ {
		if node.next == nil {
			node.next = make(map[byte]*trieCacheNodeMap[V])
		}

		node = node.next[key[i]]

		if node == nil {
			node = &trieCacheNodeMap[V]{c: key[i]}
			prev.next[key[i]] = node
		}

		prev = node
	}

	if node.value == nil {
		kv.len++
	}
	node.value = &value
}

func (kv *KVCacheMap[V]) dfs(node *trieCacheNodeMap[V], prefix []byte) ([]V, error) {
	res := []V{}
	if node.value != nil {
		res = append(res, *node.value)
	}

	if node.next != nil {
		// Need to maintain order, so we can't just iterate over map keys.
		i := byte(0)
		for {
			if node.next[i] != nil {
				next := node.next[i]
				nextRes, err := kv.dfs(next, append(prefix, next.c))
				if err != nil {
					return nil, err
				}
				res = append(res, nextRes...)
			}

			// Have to do this to avoid overflow.
			if i == 255 {
				break
			}
			i++
		}
	}

	return res, nil
}

func (kv *KVCacheMap[V]) ListByPrefix(prefix string) ([]V, error) {
	kv.mux.RLock()
	defer kv.mux.RUnlock()

	node := kv.trie
	for i := 0; i < len(prefix); i++ {
		node = node.next[prefix[i]]
		if node == nil {
			return nil, nil
		}
	}

	return kv.dfs(node, []byte(prefix))
}

// Get value by key from the underlying cache.
func (kv *KVCacheMap[V]) Get(key string) (V, error) {
	kv.mux.RLock()
	defer kv.mux.RUnlock()
	node := kv.trie
	for i := 0; i < len(key); i++ {
		if node.next[key[i]] == nil {
			return kv.zero, geche.ErrNotFound
		}
		node = node.next[key[i]]
	}

	if node.value == nil {
		return kv.zero, geche.ErrNotFound
	}

	return *node.value, nil
}

// Del key from the cache.
func (kv *KVCacheMap[V]) Del(key string) error {
	kv.mux.Lock()
	defer kv.mux.Unlock()

	node := kv.trie
	var prev *trieCacheNodeMap[V]
	for i := 0; i < len(key); i++ {
		next := node.next[key[i]]
		if next == nil {
			return nil
		}

		prev = node
		node = next
	}

	node.value = nil

	empty := true
	for k := range node.next {
		if node.next[k] != nil {
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
func (kv *KVCacheMap[V]) Snapshot() map[string]V {
	kv.mux.RLock()
	defer kv.mux.RUnlock()
	return kv.dfs2(kv.trie, []byte{})
}

// Len returns total number of elements in the underlying caches.
func (kv *KVCacheMap[V]) Len() int {
	kv.mux.RLock()
	defer kv.mux.RUnlock()
	return kv.len
}

type stackItem[V any] struct {
	node   *trieCacheNodeMap[V]
	prefix []byte
	c      byte
}

func (kv *KVCacheMap[V]) dfs2(node *trieCacheNodeMap[V], prefix []byte) map[string]V {
	res := map[string]V{}
	stack := []stackItem[V]{
		{node: node, c: 0},
	}

	if node.value != nil {
		res[string(prefix)] = *node.value
	}

	for {
		if len(stack) == 0 {
			break
		}

		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		key := append(prefix, top.prefix...)

		if top.node.value != nil {
			res[string(key)] = *top.node.value
		}

		if top.node.next[top.c] != nil {
			stack = append(stack,
				stackItem[V]{
					node:   top.node.next[top.c],
					c:      0,
					prefix: append(top.prefix, top.c),
				})
			continue
		}

		if top.c < 255 {
			stack = append(stack, stackItem[V]{node: top.node, prefix: top.prefix, c: top.c + 1})
			continue
		}
	}

	return res
}
