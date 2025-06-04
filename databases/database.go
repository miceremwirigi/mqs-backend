package databases

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartDatabase() (db *gorm.DB) {
	log.Println("Starting Database ...")
	_, db_host, db_user, db_pass, db_name, db_ssl, db_port, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Africa/Nairobi`,
		db_host, db_user, db_pass, db_name, db_port, db_ssl)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		log.Println("Successfully Started Database")
	} else {
		log.Fatal(err)
	}
	return db
}
