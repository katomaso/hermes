package storage

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"log"
)

var (
	db *pg.DB
)

func InitStorage(dbName, dbHost, dbUser, dbPass string) {
	log.Printf("Initializing PostgreSQL database %s@%s/%s\n", dbUser, dbHost, dbName)
	db = pg.Connect(&pg.Options{
		User:     dbUser,
		Password: dbPass,
		Database: dbName,
	})
}

func CreateStorage(model interface{}) {
	err := db.CreateTable(model, &orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func CloseStorage() {
	db.Close()
}

func Get(pointer interface{}, where string, params ...interface{}) error {
	return db.Model(pointer).Where(where, params...).Select()
}

func Insert(value interface{}) error {
	return db.Insert(value)
}

func Update(value interface{}) error {
	return db.Update(value)
}
