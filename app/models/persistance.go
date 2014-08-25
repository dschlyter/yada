package models

import (
    "github.com/syndtr/goleveldb/leveldb"
    "os"
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
    // TODO create timeUUID
    ret = "key"
    return
}

func get(nextKey string, limit int) (values [][]byte, err error) {
    iter := store.NewIterator(nil, nil)
    defer iter.Release()

    for iter.Next() {
        if iter.Error() != nil {
            return nil, iter.Error()
        }

        values = append(values, iter.Value())
    }

    return
}
