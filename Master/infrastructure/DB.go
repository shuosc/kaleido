package infrastructure

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

var DB *sql.DB = nil

func initDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func init() {
	DB = initDB(os.Getenv("CONNSTR"))
}
