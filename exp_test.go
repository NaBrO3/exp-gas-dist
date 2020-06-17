package main

import (
	"log"
	"math/rand"
	"os/exec"
	"testing"
)

func Benchmark_next(b *testing.B) {
	backup()
	defer restore()

	obj := &exp{}
	obj.init(0)
	defer obj.close()

	b.ResetTimer()
	defer b.StopTimer()

	for i := 0; i < b.N; i++ {
		obj.next(100000000)
	}
}

func Benchmark_claim(b *testing.B) {
	backup()
	defer restore()

	obj := &exp{}
	obj.init(0)
	defer obj.close()

	addresses := make([][]byte, b.N)
	obj.addresses(addresses)

	b.ResetTimer()
	defer b.StopTimer()

	for i := 0; i < b.N; i++ {
		obj.claim(addresses[i])
	}
}

func Benchmark_stake(b *testing.B) {
	backup()
	defer restore()

	obj := &exp{}
	obj.init(0)
	defer obj.close()

	addresses := make([][]byte, b.N)
	neos := make([]uint64, b.N)

	for i := 0; i < b.N; i++ {
		addresses[i] = randaddr()
		neos[i] = rand.Uint64() % 100000000
	}

	b.ResetTimer()
	defer b.StopTimer()

	for i := 0; i < b.N; i++ {
		obj.stake(addresses[i], neos[i])
	}
}

func Benchmark_claim_peak(b *testing.B) {
	backup()
	defer restore()

	obj := &exp{}
	obj.init(0)
	defer obj.close()

	addresses := make([][]byte, b.N)
	obj.addresses(addresses)

	for i := 0; i < pob; i++ {
		obj.next(100000000)
	}

	b.ResetTimer()
	defer b.StopTimer()

	for i := 0; i < b.N; i++ {
		obj.claim(addresses[i])
	}
}

func backup() {
	if err := exec.Command("cp", "-r", "data", "data.bak").Run(); err != nil {
		log.Fatalln(err)
	}
}

func restore() {
	if err := exec.Command("rm", "-rf", "data").Run(); err != nil {
		log.Fatalln(err)
	}
	if err := exec.Command("mv", "data.bak", "data").Run(); err != nil {
		log.Fatalln(err)
	}
}

// in worst case, one have 50-years' unclaimed GAS
const pob = 100000000
