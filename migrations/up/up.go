package main

import (
	"log"

	"github.com/miceremwirigi/mqs-backend/databases"
	"github.com/miceremwirigi/mqs-backend/models"
)

func main() {
	db := databases.StartDatabase()
	err := db.AutoMigrate(&models.Hospital{})
	if err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}

	err = db.AutoMigrate(&models.Equipment{})
	if err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}

	err = db.AutoMigrate(&models.Engineer{})
	if err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}

	err = db.AutoMigrate(&models.Service{})
	if err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}
}
