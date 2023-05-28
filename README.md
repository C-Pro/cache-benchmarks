# Benchmark geche vs other Go generic caches

I have picked some of the Go cache libraries that have the magic word "generic" in the description and put them in the same benchmark to compare.

```shell
$ go test -benchtime=10s -benchmem -bench .
goos: linux
goarch: amd64
pkg: cache_bench
cpu: Intel(R) Xeon(R) Platinum 8280 CPU @ 2.70GHz
BenchmarkEverythingParallel/MapCache-32         	64085875	       248.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/MapTTLCache-32      	58598002	       279.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/RingBuffer-32       	48229945	       315.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/ShardedMapCache-32  	234258486	        53.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/ShardedMapTTLCache-32         	231177732	        53.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/ShardedRingBuffer-32          	236979438	        48.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/Code-Hex/go-generics-cache-32         	39842918	       345.9 ns/op	       7 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/Yiling-J/theine-go-32                 	150612642	        81.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/jellydator/ttlcache-32                	29333647	       433.9 ns/op	      43 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/erni27/imcache-32                     	345577933	        35.63 ns/op	      50 B/op	       1 allocs/op
BenchmarkEverythingParallel/github.com/dgraph-io/ristretto-32                	83293519	       142.1 ns/op	      27 B/op	       1 allocs/op
BenchmarkEverythingParallel/github.com/hashicorp/golang-lru/v2-32            	35763888	       378.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/egregors/kesh-32                      	25860772	       524.1 ns/op	      84 B/op	       2 allocs/op
BenchmarkEverythingParallel/KVMapCache-32                                    	33802629	       478.4 ns/op	     109 B/op	       0 allocs/op
PASS
```
