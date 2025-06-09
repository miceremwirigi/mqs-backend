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

	err = db.Migrator().DropTable(&models.Equipment{})
	if err != nil {
		log.Fatalf("Failed to drop table: %s", err)
	}

	err = db.Migrator().DropTable(&models.Engineer{})
	if err != nil {
		log.Fatalf("Failed to drop table: %s", err)
	}

	err = db.Migrator().DropTable(&models.Service{})
	if err != nil {
		log.Fatalf("Failed to drop table: %s", err)
	}

	err = db.Migrator().DropTable(&models.User{})
	if err != nil {
		log.Fatalf("Failed to drop table: %s", err)
	}
	log.Println("Successfully dropped all tables")
}
