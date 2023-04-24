# Benchmark geche vs other Go generic caches

I have picked some of the Go cache libraries that have the magic word "generic" in the description here:

And put them in the same benchmark to compare.

```shell
$ go test -benchtime=30s -benchmem -bench .
goos: linux
goarch: amd64
pkg: cache_bench
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkEverythingParallel/MapCache-8          100000000              436.7 ns/op             8 B/op          1 allocs/op
BenchmarkEverythingParallel/MapTTLCache-8       100000000              482.8 ns/op             9 B/op          1 allocs/op
BenchmarkEverythingParallel/RingBuffer-8        123863301              277.8 ns/op             8 B/op          1 allocs/op
BenchmarkEverythingParallel/github.com/Code-Hex/go-generics-cache-8             85186693               651.4 ns/op            16 B/op          1 allocs/op
BenchmarkEverythingParallel/github.com/Yiling-J/theine-go-8                     92254014               399.8 ns/op             8 B/op          1 allocs/op
BenchmarkEverythingParallel/github.com/jellydator/ttlcache-8                    85838470               435.0 ns/op            58 B/op          1 allocs/op
BenchmarkEverythingParallel/github.com/erni27/imcache-8                         100000000              488.1 ns/op            10 B/op          1 allocs/op
PASS
ok      cache_bench     337.399s
```
