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
Benchmark_next-4                    1000              4974 ns/op             390 B/op          9 allocs/op
Benchmark_claim-4                   1000             29523 ns/op            4374 B/op         62 allocs/op
Benchmark_stake-4                   1000              8661 ns/op             541 B/op          8 allocs/op
Benchmark_claim_peak-4              1000             26968 ns/op            4402 B/op         62 allocs/op
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