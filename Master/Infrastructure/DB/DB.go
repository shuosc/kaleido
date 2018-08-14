package DB

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB = nil

func initDB(user string, password string) *sql.DB {
	connStr := "user=" + user + " password=" + password + " dbname=kaleido sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func init() {
	DB = initDB("test", "test123456")
}
