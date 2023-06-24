# Benchmark geche vs other Go generic caches

I have picked some of the Go cache libraries that have the magic word "generic" in the description and put them in the same benchmark to compare.

```shell
$ go test -benchtime=10s -benchmem -bench .
goos: linux
goarch: amd64
pkg: cache_bench
cpu: Intel(R) Xeon(R) Platinum 8358 CPU @ 2.60GHz
BenchmarkEverythingParallel/MapCache-32                 100000000              170.1 ns/op             0 B/op          0 allocs/op
BenchmarkEverythingParallel/MapTTLCache-32              90510988               198.9 ns/op             0 B/op          0 allocs/op
BenchmarkEverythingParallel/RingBuffer-32               85731428               196.8 ns/op             0 B/op          0 allocs/op
BenchmarkEverythingParallel/ShardedMapCache-32          273706551               43.51 ns/op            0 B/op          0 allocs/op
BenchmarkEverythingParallel/ShardedMapTTLCache-32               282491904               44.37 ns/op            0 B/op          0 allocs/op
BenchmarkEverythingParallel/ShardedRingBuffer-32                284756061               40.78 ns/op            0 B/op          0 allocs/op
BenchmarkEverythingParallel/github.com/Code-Hex/go-generics-cache-32            43165059               294.2 ns/op             7 B/op          0 allocs/op
BenchmarkEverythingParallel/github.com/Yiling-J/theine-go-32                    186976719               64.51 ns/op            0 B/op          0 allocs/op
BenchmarkEverythingParallel/github.com/jellydator/ttlcache-32                   29943469               376.3 ns/op            43 B/op          0 allocs/op
BenchmarkEverythingParallel/github.com/erni27/imcache-32                        531496862               23.35 ns/op           50 B/op          1 allocs/op
BenchmarkEverythingParallel/github.com/dgraph-io/ristretto-32                   100000000              108.5 ns/op            27 B/op          1 allocs/op
BenchmarkEverythingParallel/github.com/hashicorp/golang-lru/v2-32               43857675               307.1 ns/op             0 B/op          0 allocs/op
BenchmarkEverythingParallel/github.com/egregors/kesh-32                         33866130               428.7 ns/op            83 B/op          2 allocs/op
BenchmarkEverythingParallel/KVMapCache-32                                       43328151               401.2 ns/op           112 B/op          0 allocs/op
PASS
ok      cache_bench     220.443s
```
