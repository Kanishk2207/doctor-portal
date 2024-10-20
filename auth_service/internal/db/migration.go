package db

import (
	"log"
	"os"
	"strings"
)

func AutoMigrate() {
	// Read the migration file
	// migration, err := os.ReadFile("../../internal/db/migration.sql")
	migration, err := os.ReadFile("../../internal/db/migration.sql")
	if err != nil {
		log.Fatalf("Could not read migration file: %v", err)
	}

	statements := strings.Split(string(migration), ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)

		if len(statement) == 0 {
			continue
		}
		_, err := DB.Exec(statement)
		if err != nil {
			log.Fatalf("Could not execute statement: %v", err)
		}
	}

	log.Println("Database migration completed successfully.")
}
