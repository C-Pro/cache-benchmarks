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
BenchmarkEverythingParallel/sync.Map-16                        862699880        18.72 ns/op       4 B/op       0 allocs/op
BenchmarkEverythingParallel/MapCache-16                        100000000       105.3 ns/op        0 B/op       0 allocs/op
BenchmarkEverythingParallel/MapTTLCache-16                     100000000       132.5 ns/op        0 B/op       0 allocs/op
BenchmarkEverythingParallel/RingBuffer-16                      100000000       111.6 ns/op        0 B/op       0 allocs/op
BenchmarkEverythingParallel/KVCache-16                         100000000       256.6 ns/op        8 B/op       0 allocs/op
BenchmarkEverythingParallel/ShardedMapCache-16                 361509866        35.32 ns/op       0 B/op       0 allocs/op
BenchmarkEverythingParallel/ShardedMapTTLCache-16              324419283        38.84 ns/op       0 B/op       0 allocs/op
BenchmarkEverythingParallel/ShardedRingBuffer-16               360948327        32.92 ns/op       0 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/Code-Hex/go-generics-cache-16  100000000       183.2 ns/op        6 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/Yiling-J/theine-go-16          256218966        48.16 ns/op       0 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/jellydator/ttlcache-16         86884353       243.4 ns/op       43 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/erni27/imcache-16              531843435        25.58 ns/op       2 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/dgraph-io/ristretto-16         142026916        92.92 ns/op       28 B/op       1 allocs/op
BenchmarkEverythingParallel/github.com/hashicorp/golang-lru/v2-16     100000000       201.2 ns/op        0 B/op       0 allocs/op
BenchmarkEverythingParallel/github.com/egregors/kesh-16               58256961       302.9 ns/op       78 B/op       2 allocs/op
BenchmarkEverythingParallel/KVMapCache-16                             100000000       316.5 ns/op       30 B/op       0 allocs/op
```

---

## Dedicated Key-Value & Prefix Matching Benchmark

This benchmark evaluates caches that support **Prefix-based lookups** (`ListByPrefix`). We compare `KVMapCache`, `KVCache`, `github.com/hashicorp/go-immutable-radix`, and `github.com/armon/go-radix`. Sharding has been excluded since merging results across shards would disrupt sorting order without an expensive merge/sort step.

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
BenchmarkKV_Insert/KVMapCache-16                        5209327       2082 ns/op     228 B/op          3 allocs/op
BenchmarkKV_Insert/KVCache-16                          69173250        170.8 ns/op     0 B/op          0 allocs/op
BenchmarkKV_Insert/ImmutableRadix-16                    4053322       3132 ns/op    3393 B/op         50 allocs/op
BenchmarkKV_Insert/ArmonRadix-16                       46540156        251.8 ns/op    16 B/op          1 allocs/op

BenchmarkKV_Get/KVMapCache-16                         410783544        29.91 ns/op     0 B/op          0 allocs/op
BenchmarkKV_Get/KVCache-16                            444693756        27.28 ns/op     0 B/op          0 allocs/op
BenchmarkKV_Get/ImmutableRadix-16                     264676414        45.96 ns/op     0 B/op          0 allocs/op
BenchmarkKV_Get/ArmonRadix-16                         342359751        40.28 ns/op     0 B/op          0 allocs/op

BenchmarkKV_ListByPrefix/KVMapCache-16                    16365     722385 ns/op   2047707 B/op      15544 allocs/op
BenchmarkKV_ListByPrefix/KVCache-16                       61632     201852 ns/op   1299588 B/op         13 allocs/op
BenchmarkKV_ListByPrefix/ImmutableRadix-16                42810     277578 ns/op   1300180 B/op         13 allocs/op
BenchmarkKV_ListByPrefix/ArmonRadix-16                    52350     226013 ns/op   1300460 B/op         13 allocs/op
```

### Analysis & Summary

1. **Insert Performance**:
   - `KVCache` is the fastest for insert operations (`170.8 ns/op`), utilizing a mutable tree with a flat memory layout and zero allocations.
   - `ArmonRadix` is also very performant (`251.8 ns/op`), easily outperforming immutable and map-backed alternatives.
   - `ImmutableRadix` exhibits the highest insert latency (`3132 ns/op`) and allocations (`50 allocs/op`) due to path-copying/allocation costs during copy-on-write inserts.

2. **Get Performance**:
   - `KVCache` (`27.28 ns/op`) and `KVMapCache` (`29.91 ns/op`) show the best read performance, followed by `ArmonRadix` (`40.28 ns/op`) and `ImmutableRadix` (`45.96 ns/op`).

3. **Prefix Listing Performance**:
   - `KVCache` offers the fastest prefix retrieval times (`201852 ns/op`). It is backed by a flat values slice, allowing it to directly resolve node values without hash map lookups.
   - `ArmonRadix` performs a close second at `226013 ns/op`.
   - `KVMapCache` is considerably slower (`722385 ns/op`) and generates high garbage collection pressure (`15544 allocs/op`) because it must query the underlying hash map for every matching key.
