# experiment of gas distribution of neo3

In this experiment, each `exp` struct represents a gas distribution instance for each consensus/committee node and voters share the inflation GAS by their stake NEO.

# run

1. generate test data

   ```sh
    go run main.go
   ```
2. run test

   ```sh
    go test -benchmem -run=^$ -bench . -timeout=30m -benchtime=1000x
   ```

# expected outputs

```
root@3b50fac68f47:/workspaces# go test -benchmem -run=^$ -bench . -timeout=30m -benchtime=1000x
goos: linux
goarch: amd64
Benchmark_next-4                    1000              4319 ns/op             287 B/op          5 allocs/op
Benchmark_claim-4                   1000             32746 ns/op            4164 B/op         56 allocs/op
Benchmark_stake-4                   1000              9519 ns/op             537 B/op          8 allocs/op
Benchmark_claim_peak-4              1000            483743 ns/op           58947 B/op       1988 allocs/op
PASS
ok      _/workspaces/100mcal    1329.084s
```

after `100000000` blocks

```
root@3b50fac68f47:/workspaces# du -h -d 1 data
749M    data/gas
2.6G    data/blk
2.7G    data/neo
4.4M    data/l1
80K     data/l2
20K     data/l3
6.0G    data
```