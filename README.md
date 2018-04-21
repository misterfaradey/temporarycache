# temporarycache
Golang temporary cache on interfaces.


## Get the Repository

``go get github.com/misterfaradey/temporarycache``

### Benchmarks
intel core i5-7200U
```
BenchmarkWrite-4         3000000               424 ns/op              24 B/op          2 allocs/op
BenchmarkGet-4          20000000                60.4 ns/op             0 B/op          0 allocs/op
BenchmarkDeleteOld-4     5000000               463 ns/op               0 B/op          0 allocs/op
```
