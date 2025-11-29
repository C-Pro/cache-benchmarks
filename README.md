# Benchmark geche vs other Go generic caches

I have picked some of the Go cache libraries that have the magic word "generic" in the description and put them in the same benchmark to compare.

```shell
go test -benchtime=10s -benchmem -bench .
goos: linux
goarch: amd64
pkg: cache_bench
cpu: AMD Ryzen 7 PRO 8840U w/ Radeon 780M Graphics
BenchmarkEverythingParallel/MapCache-16         	100000000	       173.9 ns/op	       1 B/op	       0 allocs/op
BenchmarkEverythingParallel/MapTTLCache-16      	65010415	       382.5 ns/op	       2 B/op	       0 allocs/op
BenchmarkEverythingParallel/RingBuffer-16       	100000000	       225.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/ShardedMapCache-16  	198813898	        56.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/ShardedMapTTLCache-16         	122482419	        97.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/ShardedRingBuffer-16          	188570131	        63.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/Code-Hex/go-generics-cache-16         	55945956	       474.2 ns/op	      91 B/op	       2 allocs/op
BenchmarkEverythingParallel/github.com/Yiling-J/theine-go-16                 	71100289	       172.2 ns/op	       1 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/jellydator/ttlcache-16                	58265994	       924.8 ns/op	      30 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/erni27/imcache-16                     	236973852	        45.84 ns/op	      33 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/dgraph-io/ristretto-16                	88618468	       141.9 ns/op	      94 B/op	       2 allocs/op
BenchmarkEverythingParallel/github.com/hashicorp/golang-lru/v2-16            	78454165	       399.1 ns/op	       2 B/op	       0 allocs/op
BenchmarkEverythingParallel/github.com/egregors/kesh-16                      	68416022	       337.8 ns/op	       7 B/op	       0 allocs/op
BenchmarkEverythingParallel/KVMapCache-16                                    	 7607014	      2050 ns/op	     254 B/op	       3 allocs/op
BenchmarkEverythingParallel/KVCache-16                                       	52397652	       902.8 ns/op	      53 B/op	       0 allocs/op
BenchmarkEverythingParallel/ShardedKVCache-16                                	100000000	       150.9 ns/op	      45 B/op	       0 allocs/op
PASS
ok  	cache_bench	390.130s
```
