package utils

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

// ShouldSendReminder checks if a reminder should be sent
func ShouldSendReminder(eq Equipment) bool {
	if eq.SnoozeEmail || eq.IsDone {
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
func ReminderHTMLTemplate(equipment Equipment) string {
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
`, equipment.Name, equipment.DueDate.Format("2006-01-02"), equipment.ID, equipment.ID, equipment.ID)
}

// ReminderCronJob runs every morning and sends reminders
func ReminderCronJob(equipments []Equipment, smtpHost string, smtpPort int, smtpUser, smtpPass string, updateReminderDate func(id string, t time.Time)) {
	c := cron.New()
	c.AddFunc("0 7 * * *", func() { // Every day at 7:00 AM
		for _, eq := range equipments {
			if ShouldSendReminder(eq) {
				html := ReminderHTMLTemplate(eq)
				subject := fmt.Sprintf("Service Due: %s", eq.Name)
				// Send to engineer
				if eq.EngineerEmail != "" {
					_ = SendReminderEmail(smtpHost, smtpPort, smtpUser, smtpPass, eq.EngineerEmail, subject, html)
				}
				// Send to hospital
				if eq.HospitalEmail != "" {
					_ = SendReminderEmail(smtpHost, smtpPort, smtpUser, smtpPass, eq.HospitalEmail, subject, html)
				}
				now := time.Now()
				updateReminderDate(eq.ID, now) // <-- Save to DB
			}
		}
	})
	c.Start()
}
