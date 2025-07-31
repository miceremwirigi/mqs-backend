package utils

import (
	"log"

	"github.com/miceremwirigi/mqs-backend/models"
	"gorm.io/gorm"
)

func RunCronJobs(db *gorm.DB) {
	// Load SMTP configuration
	smtpHost, smtpPort, smtpUser, smtpPass, err := LoadSMTPConfig()
	if err != nil {
		panic(err)
	}

	// Get the list of equipments from the database
	equipments, err := models.GetEquipments(db)
	if err != nil {
		panic("Failed to retrieve equipments: " + err.Error())
	}
	
	// Schedule regular email reminder cron job
	ReminderCronJob(db, equipments, smtpHost, smtpPort, smtpUser, smtpPass, UpdateReminderDate)
	
	log.Println("Reminder cron job scheduled successfully")
}

