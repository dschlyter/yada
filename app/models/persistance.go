package models

import (
    "github.com/syndtr/goleveldb/leveldb"
    "math/rand"
    "os"
    "strconv"
)

var dbName = "level.db"
var store *leveldb.DB

func InitDB() {
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

// Crude test fixture
func SetTestMode() {
    CloseDB()
    dbName = "tests/test.db"
    os.RemoveAll(dbName) // Wipe previous data
    InitDB()
}

func save(data []byte) (err error) {
    key := []byte(createKey())
    err = store.Put(key, data, nil)
    return
}

func createKey() (ret string) {
    ret = "key" + "-" + strconv.Itoa(rand.Int())
    return
}

func get(nextKey string, limit int) (values [][]byte, err error) {
    iter := store.NewIterator(nil, nil)
    defer iter.Release()

    for iter.Next() {
        if iter.Error() != nil {
            return nil, iter.Error()
        }

        // this is slightly wasteful, maybe serialize json directly
        newBytes := make([]byte, len(iter.Value()))
        copy(newBytes, iter.Value())
        values = append(values, newBytes)
    }

    return
}
