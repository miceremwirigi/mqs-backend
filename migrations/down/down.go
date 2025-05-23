package main

import (
	"log"

	"github.com/miceremwirigi/mqs-backend/databases"
	"github.com/miceremwirigi/mqs-backend/models"
)

func main() {
	db := databases.StartDatabase()
	err := db.Migrator().DropTable(&models.Hospital{})
	if err != nil {
		log.Fatalf("Failed to drop table: %s", err)
	}
}
