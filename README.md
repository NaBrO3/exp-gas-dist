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
Benchmark_next-4            1000             13292 ns/op            3987 B/op         57 allocs/op
Benchmark_claim-4           1000            431994 ns/op           18519 B/op        213 allocs/op
Benchmark_stake-4           1000              8603 ns/op             542 B/op          8 allocs/op
PASS
ok      _/workspaces/100mcal    341.949s
```

after `100000000` blocks

```
root@3b50fac68f47:/workspaces# du -h -d 1 data
16K     data/gas
2.6G    data/blk
2.7G    data/neo
5.3G    data
```