# Benchmark geche vs other Go generic caches

I have picked some of the Go cache libraries that have the magic word "generic" in the description and put them in the same benchmark to compare.

```shell
$ go test -benchtime=30s -benchmem -bench .
goos: linux
goarch: amd64
pkg: cache_bench
cpu: Intel(R) Xeon(R) Platinum 8168 CPU @ 2.70GHz
BenchmarkEverythingParallel/MapCache-32         	223955310	       178.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/MapTTLCache-32      	187183848	       220.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/RingBuffer-32       	202385371	       205.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/ShardedMapCache-32  	610846776	        60.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/ShardedMapTTLCache-32         	632016300	        62.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/ShardedRingBuffer-32          	645748414	        55.06 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/Code-Hex/go-generics-cache-32         	131428951	       304.8 ns/op	       7 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/Yiling-J/theine-go-32                 	377735642	        90.43 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/jellydator/ttlcache-32                	90392742	       334.8 ns/op	      43 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/erni27/imcache-32                     	1000000000	        19.56 ns/op	       2 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/dgraph-io/ristretto-32                	266648240	       149.3 ns/op	      28 B/op	       1 allocs/op
BenchmarkEverythingParallel/github.com/hashicorp/golang-lru/v2-32            	100000000	       331.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/egregors/kesh-32                      	100000000	       372.7 ns/op	      46 B/op	       2 allocs/op
PASS
```
