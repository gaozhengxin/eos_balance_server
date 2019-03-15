package dao
import (
	"../config"
	"log"
	"github.com/syndtr/goleveldb/leveldb"
)

var dbPath = config.DbPath
var db *leveldb.DB

func init () {
	var err error
	db, err = leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Close(db *leveldb.DB) {
	if r := recover(); r != nil {
		log.Println(r)
	}
	db.Close()
}

func Put(key string, value string) {
	e := db.Put([]byte(key), []byte(value), nil)
	if e != nil {
		panic(e)
	}
	return
}

func Get(key string) string {
	data, e := db.Get([]byte(key), nil)
	if e != nil {
		if e.Error() == "leveldb: not found" {
			return ""
		}
		panic(e)
	}
	return string(data)
}

func Delete(key string) {
	e := db.Delete([]byte(key),nil)
	if e != nil {
		panic(e)
	}
}
