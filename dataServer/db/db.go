package db

import "github.com/syndtr/goleveldb/leveldb"

func OpenDB(dbPath string) (*leveldb.DB, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDB(db *leveldb.DB) error {
	err := db.Close()
	return err
}

func Get(db *leveldb.DB, key []byte) ([]byte, error) {
	data, err := db.Get(key, nil)
	return data, err
}

func Put(db *leveldb.DB, key []byte, value []byte) error {
	err := db.Put(key, value, nil)
	return err
}

func Del(db *leveldb.DB, key []byte) error {
	err := db.Delete(key, nil)
	return err
}
