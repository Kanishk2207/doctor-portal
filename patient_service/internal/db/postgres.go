package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dsn string) {
	var err error

	// Connect to PostgreSQL database
	DB, err = sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	log.Printf("Successfully connected to DB: %v", dsn)

	// Ping the database to verify connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}

	AutoMigrate()
}
