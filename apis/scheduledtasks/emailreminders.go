package scheduledtasks

import (
	"github.com/gofiber/fiber/v3"
	"github.com/miceremwirigi/mqs-backend/common/apis"
	"github.com/miceremwirigi/mqs-backend/models"
	"github.com/miceremwirigi/mqs-backend/utils"
)

func (h *Handler) SendEmailAllReminders(c fiber.Ctx) error {
	// This function will be used to send email reminders
	// It can be implemented later as per the requirements

	// Load SMTP configuration
	smtpHost, smtpPort, smtpUser, smtpPass, err := utils.LoadSMTPConfig()
	if err != nil {
		panic(err)
	}

	// Get the list of equipments from the database
	equipments, err := models.GetEquipments(h.DB)
	if err != nil {
		panic("Failed to retrieve equipments: " + err.Error())
	}

	go func() {
		utils.SendServiceDueRemindersImmediately(
			h.DB, equipments, smtpHost, smtpPort, smtpUser, smtpPass, utils.UpdateReminderDate)
	}()

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully sent all service due emails", nil)
}
