package models

import (
	"math/rand"
	"strconv"
	"sync/atomic"

	"github.com/syndtr/goleveldb/leveldb"
)

var store *leveldb.DB
var suffix uint64

const (
	DB      = "level.db"
	DB_TEST = "tests/test.db"
)

// These methods are not goroutine safe, use in init and testing only
func InitDB(dbName string) {
	CloseDB()

	db, err := leveldb.OpenFile(dbName, nil)
	if err != nil {
		panic(err)
	}
	store = db
	suffix = uint64(rand.Int63())
}

func CloseDB() {
	if store != nil {
		err := store.Close()
		if err != nil {
			panic(err)
		}
		store = nil
	}
}

func save(keyPrefix string, data []byte) (err error) {
	key := []byte(createKey(keyPrefix))
	err = store.Put(key, data, nil)
	return
}

func createKey(keyPrefix string) (ret string) {
	ctr := atomic.AddUint64(&suffix, 1)
	ret = keyPrefix + "-" + strconv.FormatUint(ctr, 10)
	return
}

type KV struct {
	Key   string
	Value []byte
}

func getBlobs() (values []KV, err error) {
	iter := store.NewIterator(nil, nil)
	defer iter.Release()

	keyExists := iter.Last()
	for keyExists {
		if iter.Error() != nil {
			return nil, iter.Error()
		}

		// this is slightly wasteful, maybe serialize json directly
		newBytes := make([]byte, len(iter.Value()))
		copy(newBytes, iter.Value())
		values = append(values, KV{string(iter.Key()), newBytes})

		keyExists = iter.Prev()
	}

	return
}
