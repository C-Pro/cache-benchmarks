# Benchmark geche vs other Go generic caches

I have picked some of the Go cache libraries that have the magic word "generic" in the description and put them in the same benchmark to compare.

```shell
$ go test -benchtime=30s -benchmem -bench .
goos: darwin
goarch: arm64
pkg: cache_bench
BenchmarkEverythingParallel/MapCache-10                 332130052              133.3 ns/op             0 B/op          0 allocs/op
BenchmarkEverythingParallel/MapTTLCache-10              234690624              205.4 ns/op             0 B/op          0 allocs/op
BenchmarkEverythingParallel/RingBuffer-10               441694302               86.82 ns/op            0 B/op          0 allocs/op
BenchmarkEverythingParallel/github.com/Code-Hex/go-generics-cache-10            191366336              198.8 ns/op             7 B/op          0 allocs/op
BenchmarkEverythingParallel/github.com/Yiling-J/theine-go-10                    367538067              100.7 ns/op             0 B/op          0 allocs/op
BenchmarkEverythingParallel/github.com/jellydator/ttlcache-10                   136785907              262.4 ns/op            43 B/op          0 allocs/op
BenchmarkEverythingParallel/github.com/erni27/imcache-10                        226084180              179.2 ns/op             2 B/op          0 allocs/op
BenchmarkEverythingParallel/github.com/dgraph-io/ristretto-10                   466729495               80.03 ns/op           30 B/op          1 allocs/op
BenchmarkEverythingParallel/github.com/hashicorp/golang-lru/v2-10               193697901              216.5 ns/op             0 B/op          0 allocs/op
PASS
ok      cache_bench     496.390s
```
