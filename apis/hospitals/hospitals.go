package hospitals

import (
	"github.com/gofiber/fiber/v3"
	"github.com/miceremwirigi/mqs-backend/common/apis"
	"github.com/miceremwirigi/mqs-backend/models"
)

func (h *Handler) GetAllHospitals(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	var hospitals []models.Equipment
	err := tx.Find(&hospitals).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error getting hospitals", nil)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", nil)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived hositals", hospitals)
}
