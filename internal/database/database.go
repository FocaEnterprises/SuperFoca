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

	// The connection is configured by libpq using Postgres environment variables.
	// These environment variables are listed here: https://www.postgresql.org/docs/current/libpq-envars.html
	db, err = sql.Open("postgres", "")

	if err != nil {
		log.Fatalf("failed to establish database connection: %s", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to reach database: %s", err)
	}

	log.Println("succesfully established database connection")
}
