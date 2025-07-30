package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/mailersend/mailersend-go"
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
	return time.Since(*eq.LastReminderDate) > 30*time.Second
	// return time.Since(*eq.LastReminderDate) > 5*24*time.Hour
}

// SendReminderEmail sends an email using SMTP
func SendReminderEmail(smtpHost string, smtpPort int, smtpUser, smtpPass, to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: smtpHost} // Explicit TLS config is good practice

	return d.DialAndSend(m)
}

func HTTPEmailSenderAlternative(from, to, subject, htmlBody string) error {
	mailersendAPIKey := os.Getenv("mailer_send_api_key")
	if mailersendAPIKey == "" {
		return fmt.Errorf("MAILERSEND_API_KEY environment variable is not set")
	}
	// Create an instance of the mailersend client
	ms := mailersend.NewMailersend(mailersendAPIKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fromStruct := mailersend.From{
		Name:  "MQS medical equipment service",
		Email: from,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  "Dear Client",
			Email: to,
		},
	}

	// Send in 5 seconds
	sendAt := time.Now().Add(time.Second * 5).Unix()

	tags := []string{"important", "service-reminder", "equipment"}

	message := ms.Email.NewMessage()

	message.SetFrom(fromStruct)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(htmlBody)
	message.SetTags(tags)
	message.SetSendAt(sendAt)
	message.SetInReplyTo("client-id")

	res, _ := ms.Email.Send(ctx, message)

	log.Print(res.Header.Get("x-message-id"))
	return nil
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
	var reminderSent = false
	nrb, _ := time.LoadLocation("Africa/Nairobi")
	c := cron.New(cron.WithLocation(nrb))
	// c.AddFunc("@every 30s", func() { // Every 30 seconds for testing
	c.AddFunc("0 7 * * *", func() { // Every day at 7:00 AM

		for _, eq := range equipments {
			if ShouldSendReminder(eq) {
				engineersEmail := eq.EngineersEmail()
				html := ReminderHTMLTemplate(eq)
				subject := fmt.Sprintf("Service Due: %s", eq.Name)
				// Send to engineer
				if engineersEmail != "" {
					err := SendReminderEmail(smtpHost, smtpPort, smtpUser, smtpPass, engineersEmail, subject, html)
					if err != nil {
						if opErr, ok := err.(*net.OpError); ok {
							if opErr.Op == "write" && opErr.Err.Error() == "broken pipe" {
								log.Printf("Broken pipe error sending reminder email to engineer %s: %v", engineersEmail, err)
							}
							if opErr.Timeout() {
								log.Printf("Timeout sending reminder email to engineer %s: %v", engineersEmail, err)
							}
						}
						if opErr, ok := err.(*net.OpError); ok && opErr.Temporary() {
							log.Printf("Temporary error sending reminder email to engineer %s: %v", engineersEmail, err)
						} else {
							log.Printf("Failed to send reminder email to engineer %s: %v", engineersEmail, err)
						}						
						HTTPEmailSenderAlternative(smtpUser, engineersEmail, subject, html)
					} else if !reminderSent {
						reminderSent = true
					}
				}
				// Send to hospital
				if eq.Hospital.Email != "" {
					err := SendReminderEmail(smtpHost, smtpPort, smtpUser, smtpPass, eq.Hospital.Email, subject, html)
					if err != nil {
						if opErr, ok := err.(*net.OpError); ok {
							if opErr.Op == "write" && opErr.Err.Error() == "broken pipe" {
								log.Printf("Broken pipe error sending reminder email to hospital%s: %v", eq.Hospital.Email, err)
							}
							if opErr.Timeout() {
								log.Printf("Timeout sending reminder email to hospital%s: %v", eq.Hospital.Email, err)
							}
						}
						if opErr, ok := err.(*net.OpError); ok && opErr.Temporary() {
							log.Printf("Temporary error sending reminder email to hospital%s: %v", eq.Hospital.Email, err)
						} else {
							log.Printf("Failed to send reminder email to hospital%s: %v", eq.Hospital.Email, err)
						}						
						HTTPEmailSenderAlternative(smtpUser, eq.Hospital.Email, subject, html)
					} else if !reminderSent {
						reminderSent = true
					}
				}
				now := time.Now()
				updateReminderDate(db, eq.ID.String(), now) // <-- Save to DB
			}
		}
		if reminderSent {
			log.Println("Reminder emails sent successfully")
			reminderSent = false // Reset for the next run
		} else {
			log.Println("No reminders to send at this time")
		}

	})
	c.Start()
}

func SendServiceDueRemindersImmediately(db *gorm.DB, equipments []models.Equipment, smtpHost string, smtpPort int, smtpUser, smtpPass string, updateReminderDate func(db *gorm.DB, id string, t time.Time)) {
	var reminderSent = false
	for _, eq := range equipments {
		if ShouldSendReminder(eq) {
			engineersEmail := eq.EngineersEmail()
			html := ReminderHTMLTemplate(eq)
			subject := fmt.Sprintf("Service Due: %s", eq.Name)
			// Send to engineer
			if engineersEmail != "" {
				err := SendReminderEmail(smtpHost, smtpPort, smtpUser, smtpPass, engineersEmail, subject, html)
				if err != nil {
					if opErr, ok := err.(*net.OpError); ok {
						if opErr.Op == "write" && opErr.Err.Error() == "broken pipe" {
							log.Printf("Broken pipe error sending reminder email to engineer %s: %v", engineersEmail, err)
						}
						if opErr.Timeout() {
							log.Printf("Timeout sending reminder email to engineer %s: %v", engineersEmail, err)
						}
					}
					if opErr, ok := err.(*net.OpError); ok && opErr.Temporary() {
						log.Printf("Temporary error sending reminder email to engineer %s: %v", engineersEmail, err)
					} else {
						log.Printf("Failed to send reminder email to engineer %s: %v", engineersEmail, err)
					}
				} else if !reminderSent {
					reminderSent = true
				}
			}
			// Send to hospital
			if eq.Hospital.Email != "" {
				err := SendReminderEmail(smtpHost, smtpPort, smtpUser, smtpPass, eq.Hospital.Email, subject, html)
				if err != nil {
					log.Printf("Failed to send reminder email to hospital %s: %v", eq.Hospital.Email, err)
				} else if !reminderSent {
					reminderSent = true
				}
			}
			now := time.Now()
			updateReminderDate(db, eq.ID.String(), now) // <-- Save to DB
		}
	}
	if reminderSent {
		log.Println("Reminder emails sent successfully")
		reminderSent = false // Reset for the next run
	} else {
		log.Println("No reminders to send at this time")
	}
}
