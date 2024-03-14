package sql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func NewConn() *sql.DB { // need to write up some config file read
	connStr := "user=server password=server dbname=userdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
