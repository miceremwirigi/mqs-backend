package databases

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartDatabase() (db *gorm.DB){
	log.Println("Starting Database ...")
	db_name, db_pass, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	dsn := fmt.Sprintf(`host=localhost user=mqs_user
		password=%s dbname=%s port=5432 sslmode= disable TimeZone=Africa/Nairobi`, db_pass, db_name)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		log.Println("Successfully Started Database")

	} else {
		log.Fatal(err)
	}
	return db
}
