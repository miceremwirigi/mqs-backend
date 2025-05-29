package hospitals

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/miceremwirigi/mqs-backend/common"
	"github.com/miceremwirigi/mqs-backend/common/apis"
	"github.com/miceremwirigi/mqs-backend/models"
)

func (h *Handler) GetAllHospitals(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	var hospitals []models.Hospital
	err := tx.Find(&hospitals).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error getting hospitals", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived hositals", hospitals)
}

func (h *Handler) GetHospital(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var hospital models.Hospital
	err := tx.Where("id = ?", id).Find(&hospital).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving hospital", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived hospital datail", hospital)
}

func (h *Handler) GetHospitalHtml(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var hospital models.Hospital
	err := tx.Where("id = ?", id).Find(&hospital).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving hospital", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	hospitalMap, err := common.StructToMap(hospital)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error converting hospital struct to map", err)
	}

	responseFiberMap := fiber.Map(hospitalMap)

	return c.Render("hospital", responseFiberMap)
}

func (h *Handler) AddHospital(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	hospital := &models.Hospital{}

	err := json.Unmarshal(c.Body(), hospital)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", err.Error())
	}

	if hospital.Name == "" {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", errors.New("error: empty hospital name"))
	}

	err = tx.Create(hospital).Error
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

	return apis.GeneralApiResponse(c, apis.StatusCreatedResponseCode, "successfully added hospital", &hospital)
}

func (h *Handler) UpdateHospital(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	id := c.Params("id")

	hospital := &models.Hospital{}
	err := json.Unmarshal(c.Body(), hospital)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", err.Error())
	}

	if hospital.Name == "" {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", errors.New("error: empty hospital name"))
	}
	err = tx.Where("id = ?", id).Updates(hospital).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction on hospital update", err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	log.Println("Updating Hospital ...")

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully uptdated hospital", &hospital)
}

func (h *Handler) DeleteHospital(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	id := c.Params("id")

	// Step 1: Fetch all equipments belonging to the hospital
	var equipments []models.Equipment
	if err := tx.Where("hospital_id = ?", id).Find(&equipments).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error fetching equipments for hospital", err.Error())
	}

	// Step 2: For each equipment, remove associations in serviced_equipments and delete the equipment
	for _, equipment := range equipments {
		if err := tx.Exec("DELETE FROM serviced_equipments WHERE equipment_id = ?", equipment.ID).Error; err != nil {
			tx.Rollback()
			return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
				"error deleting related serviced_equipments", err.Error())
		}
		if err := tx.Where("id = ?", equipment.ID).Delete(&models.Equipment{}).Error; err != nil {
			tx.Rollback()
			return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
				"error deleting equipment", err.Error())
		}
	}

	// Step 3: Delete the hospital
	hospital := &models.Hospital{}
	if err := tx.Where("id = ?", id).Delete(hospital).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction on hospital delete", err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully deleted hospital and related equipments", nil)
}
