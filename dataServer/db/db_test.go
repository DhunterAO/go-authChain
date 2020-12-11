package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"strings"
	"testing"
)

func TestDB(t *testing.T) {
	// test open and close database
	db, err := OpenDB("../data/leveldb")
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		err := CloseDB(db)
		if err != nil {
			t.Error(err)
		}
	}()

	// test add and query key-value pair in db
	key1 := "hello"
	value1 := "world"
	err = Put(db, []byte(key1), []byte(value1))
	if err != nil {
		t.Error(err)
		return
	}

	data, err := Get(db, []byte(key1))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(data[:]))
	if strings.Compare(string(data[:]), value1) != 0 {
		t.Errorf("Not Equal: " + string(data[:]) + " vs " + value1)
		return
	}

	// test update value for key in db
	value2 := "world2"
	err = Put(db, []byte(key1), []byte(value2))
	if err != nil {
		t.Error(err)
		return
	}

	data, err = Get(db, []byte(key1))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(data[:]))
	if strings.Compare(string(data[:]), value2) != 0 {
		t.Errorf("Not Equal: " + string(data[:]) + " vs " + value1)
		return
	}

	// test del key in db
	err = Del(db, []byte(key1))
	if err != nil {
		t.Error(err)
		return
	}

	data, err = Get(db, []byte(key1))
	if err != leveldb.ErrNotFound {
		t.Error(err)
		return
	}
}
