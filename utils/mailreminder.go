package utils

import (
	"fmt"
	"time"

	"github.com/miceremwirigi/mqs-backend/models"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

// ShouldSendReminder checks if a reminder should be sent
func ShouldSendReminder(eq models.Equipment) bool {
	isDone := eq.IsDone()
	if eq.SnoozeEmail || isDone {
		return false
	}
	if eq.LastReminderDate == nil {
		return true
	}
	return time.Since(*eq.LastReminderDate) > 5*24*time.Hour
}

// SendReminderEmail sends an email using SMTP
func SendReminderEmail(smtpHost string, smtpPort int, smtpUser, smtpPass, to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	return d.DialAndSend(m)
}

// ReminderHTMLTemplate returns the HTML email body
func ReminderHTMLTemplate(equipment models.Equipment) string {
	dueDate := equipment.DueDate()
	return fmt.Sprintf(`
<html>
<body>
  <h2>Equipment Service Reminder</h2>
  <p>This is a reminder that the equipment <b>%s</b> is due for service on <b>%s</b>.</p>
  <p>
    <a href="https://yourapp.example.com/set-done?id=%s">Set as Done</a> |
    <a href="https://yourapp.example.com/do-not-remind?id=%s">Do Not Remind Me</a> |
    <a href="https://yourapp.example.com/snooze?id=%s">Snooze 5 Days</a>
  </p>
  <p>Thank you.</p>
</body>
</html>
`, equipment.Name, dueDate.Format("2006-01-02"), equipment.ID, equipment.ID, equipment.ID)
}

// Updates the last date a service reminder for an equipment was sent in the database.
func UpdateReminderDate(db *gorm.DB, id string, t time.Time) {
	db.Model(&models.Equipment{}).Where("id = ?", id).Update("last_reminder_date", t)
}

// ReminderCronJob runs every morning and sends equipment service due reminders
func ReminderCronJob(db *gorm.DB, equipments []models.Equipment, smtpHost string, smtpPort int, smtpUser, smtpPass string, updateReminderDate func(db *gorm.DB, id string, t time.Time)) {
	c := cron.New()
	// c.AddFunc("@every 30s", func() { // Every 30 seconds for testing
	c.AddFunc("0 7 * * *", func() { // Every day at 7:00 AM
		for _, eq := range equipments {
			if ShouldSendReminder(eq) {
				engineersEmail := eq.EngineersEmail()
				html := ReminderHTMLTemplate(eq)
				subject := fmt.Sprintf("Service Due: %s", eq.Name)
				// Send to engineer
				if engineersEmail != "" {
					_ = SendReminderEmail(smtpHost, smtpPort, smtpUser, smtpPass, engineersEmail, subject, html)
				}
				// Send to hospital
				if eq.Hospital.Email != "" {
					_ = SendReminderEmail(smtpHost, smtpPort, smtpUser, smtpPass, eq.Hospital.Email, subject, html)
				}
				now := time.Now()
				updateReminderDate(db, eq.ID.String(), now) // <-- Save to DB
			}
		}
	})
	c.Start()
}
