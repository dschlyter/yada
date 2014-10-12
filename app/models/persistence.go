package models

import (
    "github.com/syndtr/goleveldb/leveldb"
    "math/rand"
    "strconv"
)

var store *leveldb.DB

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
    ret = keyPrefix + "-" + strconv.Itoa(rand.Int())
    return
}

// TODO actually care about limit on persistance side?
func get(nextKey string, limit int) (values [][]byte, err error) {
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
        values = append(values, newBytes)

        keyExists = iter.Prev()
    }

    return
}
