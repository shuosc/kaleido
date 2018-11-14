package DB

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)

var DB *sql.DB = nil

func initDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	for i := 0; i < 5; i++ {
		if err != nil {
			fmt.Println("Cannot connect to db, retrying...")
			time.Sleep(10)
			db, err = sql.Open("postgres", connStr)
		} else {
			break
		}
	}
	if err != nil {
		panic(err)
	}
	return db
}

func init() {
	DB = initDB(os.Getenv("CONNSTR"))
}
