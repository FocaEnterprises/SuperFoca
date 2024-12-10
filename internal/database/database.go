package database

import (
	"database/sql"
	"log"
)

var (
	db *sql.DB
)

func Init() {
	var err error

	db, err = sql.Open("postgres", "")

	if err != nil {
		log.Fatalf("failed to establish database connection: %s", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to reach database: %s", err)
	}

	log.Println("succesfully established database connection")
}
