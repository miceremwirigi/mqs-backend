package scheduledtasks

import (
	"github.com/gofiber/fiber/v3"
	"github.com/miceremwirigi/mqs-backend/common/apis"
	"github.com/miceremwirigi/mqs-backend/utils"
)

func (h *Handler) SendEmailAllReminders(c fiber.Ctx) error {
	// This function will be used to send email reminders
	// It can be implemented later as per the requirements

	// go func() {
		utils.RunCronJobs(h.DB)
	// }()

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully sent all service due emails", nil)
}
