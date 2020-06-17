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
func (me *exp) init() {
	me.dbNEO = db("data/neo")
	me.dbBLK = db("data/blk")
	me.dbGAS = db("data/gas")

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
		nBLK := uint64(0)
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
func (me *exp) next(block uint64, inflation uint64) {
	// nGAS: calculate distributed GAS per NEO in current block
	nGAS := inflation / me.count
	bGAS := make([]byte, 8)
	bBLK := make([]byte, 8)

	nGAS += me.getgas(block)
	binary.BigEndian.PutUint64(bGAS, nGAS)
	binary.BigEndian.PutUint64(bBLK, block)

	// write nGAS to db
	if err := me.dbGAS.Put(bBLK, bGAS, nil); err != nil {
		log.Fatalln(err)
	}
}

// quit staking and claim GAS
func (me *exp) claim(block uint64, address []byte) uint64 {
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
	nBLK := binary.BigEndian.Uint64(bBLK)

	gas := me.getgas(block)
	gas -= me.getgas(nBLK)

	if err := me.dbNEO.Delete(address, nil); err != nil {
		log.Fatalln(err)
	}
	if err := me.dbBLK.Delete(address, nil); err != nil {
		log.Fatalln(err)
	}
	return gas * nNEO
}

// stake
func (me *exp) stake(block uint64, address []byte, nNEO uint64) {
	bNEO := make([]byte, 8)
	bBLK := make([]byte, 8)
	binary.BigEndian.PutUint64(bNEO, nNEO)
	binary.BigEndian.PutUint64(bBLK, block)

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

// get data from dbGAS
func (me *exp) getgas(block uint64) uint64 {
	bBLK := make([]byte, 8)
	binary.BigEndian.PutUint64(bBLK, block)
	iter := me.dbNEO.NewIterator(nil, nil)
	if iter.Seek(bBLK) == false {
		if iter.Prev() == false {
			return 0
		}
	}
	bGAS := iter.Value()
	return binary.BigEndian.Uint64(bGAS)
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
