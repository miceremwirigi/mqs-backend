package engineers

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/miceremwirigi/mqs-backend/common"
	"github.com/miceremwirigi/mqs-backend/common/apis"
	"github.com/miceremwirigi/mqs-backend/models"
)

func (h *Handler) GetAllEngineers(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	var engineers []models.Engineer
	err := tx.Find(&engineers).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error getting engineers", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived engineers", engineers)
}

func (h *Handler) GetEngineer(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var engineer models.Engineer
	err := tx.Where("id = ?", id).Find(&engineer).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving engineer", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived engineer datail", engineer)
}

func (h *Handler) GetEngineerHtml(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var engineer models.Engineer
	err := tx.Where("id = ?", id).Find(&engineer).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving engineer", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	engineerMap, err := common.StructToMap(engineer)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error converting engineer struct to map", err)
	}

	responseFiberMap := fiber.Map(engineerMap)

	return c.Render("engineer", responseFiberMap)
}

func (h *Handler) AddEngineer(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	engineer := &models.Engineer{}

	err := json.Unmarshal(c.Body(), engineer)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", err.Error())
	}

	if engineer.Name == "" {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", errors.New("error: empty engineer name"))
	}

	err = tx.Create(engineer).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction", err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	return apis.GeneralApiResponse(c, apis.StatusCreatedResponseCode, "successfully added engineer", &engineer)
}

func (h *Handler) UpdateEngineer(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	id := c.Params("id")

	engineer := &models.Engineer{}
	err := json.Unmarshal(c.Body(), engineer)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", err.Error())
	}

	if engineer.Name == "" {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", errors.New("error: empty engineer name"))
	}
	err = tx.Where("id = ?", id).Updates(engineer).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction on engineer update", err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	log.Println("Updating Engineer ...")

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully uptdated engineer", &engineer)
}

func (h *Handler) DeleteEngineer(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	id := c.Params("id")

	// Remove associations in serviced_by join table before deleting engineer
	if err := tx.Exec("DELETE FROM serviced_by WHERE engineer_id = ?", id).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error deleting related serviced_by entries", err.Error())
	}

	engineer := &models.Engineer{}
	if err := tx.Where("id = ?", id).Delete(engineer).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction on engineer delete", err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully deleted engineer", nil)
}
