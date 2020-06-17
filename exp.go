package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

// GAS distribution instance
type exp struct {
	// KV {ADDRESS -> NEO BALANCE}
	dbNEO *leveldb.DB
	// KV {ADDRESS -> BLOCK INDEX WHEN START STAKING}
	dbBLK *leveldb.DB
	// KV {BLOCK INDEX -> DISTRIBUTED GAS PER NEO}
	dbGAS *leveldb.DB
	// current block index
	block uint64
	// total vote count
	count uint64
}

// remove data
func (me *exp) reset() {
	if err := os.RemoveAll("data/"); err != nil {
		log.Fatalln(err)
	}
}

// load leveldb and init
func (me *exp) init(block uint64) {
	me.dbNEO = db("data/neo")
	me.dbBLK = db("data/blk")
	me.dbGAS = db("data/gas")

	me.block = block

	// vote counting
	iter := me.dbNEO.NewIterator(nil, nil)
	for iter.Next() {
		bNEO := iter.Value()
		me.count += binary.BigEndian.Uint64(bNEO)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		log.Fatalln(err)
	}
}

// generate testing data
func (me *exp) gen(size int) {
	for i := 0; i < size; i++ {
		address := randaddr()

		nNEO := rand.Uint64()%255 + 1
		nBLK := me.block
		bNEO := make([]byte, 8)
		bBLK := make([]byte, 8)

		binary.BigEndian.PutUint64(bNEO, nNEO)
		binary.BigEndian.PutUint64(bBLK, nBLK)

		me.count += nNEO

		if err := me.dbNEO.Put(address, bNEO, nil); err != nil {
			log.Fatalln(err)
		}
		if err := me.dbBLK.Put(address, bBLK, nil); err != nil {
			log.Fatalln(err)
		}
	}

}

// next block
func (me *exp) next(inflation uint64) {
	me.block++
	// nGAS: calculate distributed GAS per NEO in current block
	nGAS := inflation / me.count
	nBLK := me.block
	bGAS := make([]byte, 8)
	bBLK := make([]byte, 8)

	nCUR := me.block - 1
	bCUR := make([]byte, 8)
	binary.BigEndian.PutUint64(bCUR, nCUR)

	if bTMP, err := me.dbGAS.Get(bCUR, nil); err == nil {
		nGAS += binary.BigEndian.Uint64(bTMP)
	}

	binary.BigEndian.PutUint64(bGAS, nGAS)
	binary.BigEndian.PutUint64(bBLK, nBLK)

	// write nGAS to db
	if err := me.dbGAS.Put(bBLK, bGAS, nil); err != nil {
		log.Fatalln(err)
	}
}

// quit staking and claim GAS
func (me *exp) claim(address []byte) uint64 {
	bNEO, err := me.dbNEO.Get(address, nil)
	if err == leveldb.ErrNotFound {
		return 0
	}
	if err != nil {
		log.Fatalln(err)
	}
	bBLK, err := me.dbBLK.Get(address, nil)
	if err != nil {
		log.Fatalln(err)
	}
	nNEO := binary.BigEndian.Uint64(bNEO)
	bCUR := make([]byte, 8)
	binary.BigEndian.PutUint64(bCUR, me.block)

	var gas uint64
	if bTMP, err := me.dbGAS.Get(bCUR, nil); err == nil {
		gas += binary.BigEndian.Uint64(bTMP)
	}

	if bTMP, err := me.dbGAS.Get(bBLK, nil); err == nil {
		gas -= binary.BigEndian.Uint64(bTMP)
	}

	if err := me.dbNEO.Delete(address, nil); err != nil {
		log.Fatalln(err)
	}
	if err := me.dbBLK.Delete(address, nil); err != nil {
		log.Fatalln(err)
	}
	return gas * nNEO
}

// stake
func (me *exp) stake(address []byte, nNEO uint64) {
	nBLK := me.block
	bNEO := make([]byte, 8)
	bBLK := make([]byte, 8)
	binary.BigEndian.PutUint64(bNEO, nNEO)
	binary.BigEndian.PutUint64(bBLK, nBLK)

	if err := me.dbNEO.Put(address, bNEO, nil); err != nil {
		log.Fatalln(err)
	}
	if err := me.dbNEO.Put(address, bBLK, nil); err != nil {
		log.Fatalln(err)
	}
}

// dump current voting address to `ret`
func (me *exp) addresses(ret [][]byte) {
	iter := me.dbNEO.NewIterator(nil, nil)
	for i := range ret {
		if iter.Next() == false {
			break
		}
		key := iter.Key()
		ret[i] = bytes.Repeat(key, 1)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		log.Fatalln(err)
	}
	return
}

// close instance
func (me *exp) close() {
	defer me.dbNEO.Close()
	defer me.dbBLK.Close()
	defer me.dbGAS.Close()
}

// open db
func db(path string) *leveldb.DB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

// generate random address
func randaddr() []byte {
	address := make([]byte, 20)
	for i := range address {
		address[i] = byte(rand.Intn(0x100))
	}
	return address
}
