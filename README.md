# Benchmark geche vs other Go generic caches

I have picked some of the Go cache libraries that have the magic word "generic" in the description and put them in the same benchmark to compare.

## General Cache Benchmark (Everything Parallel)

Run with:
```shell
$ go test -benchtime=10s -benchmem -bench BenchmarkEverythingParallel
```

### Results
```
goos: linux
goarch: amd64
pkg: cache_bench
cpu: AMD Ryzen 7 PRO 8840U w/ Radeon 780M Graphics
BenchmarkEverythingParallel/sync.Map-16                        910052367        18.79 ns/op       4 B/op       0 allocs/op
BenchmarkEverythingParallel/MapCache-16                        127467838       108.8 ns/op        0 B/op       0 allocs/op
BenchmarkEverythingParallel/MapTTLCache-16                     100000000       135.5 ns/op        0 B/op       0 allocs/op
BenchmarkEverythingParallel/RingBuffer-16                      100000000       101.3 ns/op        0 B/op       0 allocs/op
BenchmarkEverythingParallel/KVCache-16                         100000000       241.4 ns/op        7 B/op       0 allocs/op
BenchmarkEverythingParallel/ShardedMapCache-16                 365753448        34.78 ns/op       0 B/op       0 allocs/op
BenchmarkEverythingParallel/ShardedMapTTLCache-16              331130448        38.13 ns/op       0 B/op       0 allocs/op
BenchmarkEverythingParallel/ShardedRingBuffer-16               373664306        31.94 ns/op       0 B/op       0 allocs/op
BenchmarkEverythingParallel/ShardedKVCache-16                  225799561        56.68 ns/op       5 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/Code-Hex/go-generics-cache-16  100000000       168.1 ns/op        6 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/Yiling-J/theine-go-16          270983552        43.57 ns/op       0 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/jellydator/ttlcache-16         94445258       236.0 ns/op       43 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/erni27/imcache-16              537931734        24.16 ns/op       2 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/dgraph-io/ristretto-16         138824774        88.53 ns/op       28 B/op       1 allocs/op
BenchmarkEverythingParallel/github.com/hashicorp/golang-lru/v2-16     100000000       195.5 ns/op        0 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/egregors/kesh-16               63242784       308.6 ns/op       77 B/op       2 allocs/op
BenchmarkEverythingParallel/KVMapCache-16                             100000000       311.8 ns/op       29 B/op       0 allocs/op
```

---

## Dedicated Key-Value & Prefix Matching Benchmark

This benchmark evaluates caches that support **Prefix-based lookups** (`ListByPrefix`). We compare `KVMapCache`, `KVCache` (both in bare and sharded versions), and `github.com/hashicorp/go-immutable-radix`.

Run with:
```shell
$ go test -benchtime=10s -benchmem -run=^$ -bench=BenchmarkKV_
```

### Results
```
goos: linux
goarch: amd64
pkg: cache_bench
cpu: AMD Ryzen 7 PRO 8840U w/ Radeon 780M Graphics
BenchmarkKV_Insert/KVMapCache-16                        5166356       2073 ns/op
BenchmarkKV_Insert/KVMapCacheSharded-16                 4846998       2105 ns/op
BenchmarkKV_Insert/KVCache-16                          65999010        165.4 ns/op
BenchmarkKV_Insert/KVCacheSharded-16                  431787915         28.79 ns/op
BenchmarkKV_Insert/ImmutableRadix-16                    4028560       3075 ns/op
BenchmarkKV_Insert/ImmutableRadixSharded-16            16623939        740.5 ns/op

BenchmarkKV_Get/KVMapCache-16                         399924712         29.68 ns/op
BenchmarkKV_Get/KVMapCacheSharded-16                 1000000000          8.111 ns/op
BenchmarkKV_Get/KVCache-16                            445533186         35.51 ns/op
BenchmarkKV_Get/KVCacheSharded-16                     653723193         18.56 ns/op
BenchmarkKV_Get/ImmutableRadix-16                     277490256         42.85 ns/op
BenchmarkKV_Get/ImmutableRadixSharded-16              352599308         34.35 ns/op

BenchmarkKV_ListByPrefix/KVMapCache-16                    17450     685353 ns/op
BenchmarkKV_ListByPrefix/KVMapCacheSharded-16             20186     593949 ns/op
BenchmarkKV_ListByPrefix/KVCache-16                       61506     193117 ns/op
BenchmarkKV_ListByPrefix/KVCacheSharded-16                44036     272090 ns/op
BenchmarkKV_ListByPrefix/ImmutableRadix-16                44604     268574 ns/op
BenchmarkKV_ListByPrefix/ImmutableRadixSharded-16         33608     353625 ns/op
```

### Analysis & Summary

1. **Insert Performance**:
   - `KVCache` dominates here, particularly when sharded (`KVCacheSharded` at `28.79 ns/op`), because it uses a mutable tree and a flat memory layout.
   - `ImmutableRadix` is much slower for inserts due to copy-on-write overhead (`3075 ns/op`).

2. **Get Performance**:
   - `KVMapCacheSharded` has the lowest latency (`8.111 ns/op`) because it is backed by a sharded flat hash map.
   - `KVCacheSharded` is also extremely fast (`18.56 ns/op`).

3. **Prefix Listing Performance**:
   - `KVCache` (bare) offers the fastest prefix retrieval times (`193117 ns/op`). It is backed by a flat values slice, allowing it to directly resolve node values without hash map lookups.
   - `ImmutableRadix` (bare) is a close second (`268574 ns/op`).
   - Sharding generally increases prefix matching latencies slightly because prefix queries must traverse and merge findings from all individual shards sequentially.
