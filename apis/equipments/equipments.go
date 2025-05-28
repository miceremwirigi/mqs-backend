package equipments

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/miceremwirigi/mqs-backend/common"
	"github.com/miceremwirigi/mqs-backend/common/apis"
	"github.com/miceremwirigi/mqs-backend/models"
)

func (h *Handler) GetAllEquipments(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	var equipments []models.Equipment
	err := tx.Preload("Hospital").Find(&equipments).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error getting equipments", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived equipments", equipments)
}

func (h *Handler) GetEquipment(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var equipment models.Equipment
	err := tx.Preload("Hospital").Where("id = ?", id).Find(&equipment).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving equipment", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived equipment datail", equipment)
}

func (h *Handler) GetEquipmentHtml(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var equipment models.Equipment
	err := tx.Preload("Hospital").Where("id = ?", id).Find(&equipment).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving equipment", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	equipmentMap, err := common.StructToMap(equipment)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error converting equipment struct to map", err)
	}

	responseFiberMap := fiber.Map(equipmentMap)

	return c.Render("equipment", responseFiberMap)
}

func (h *Handler) AddEquipment(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	equipment := models.Equipment{}

	err := json.Unmarshal(c.Body(), &equipment)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", err.Error())
	}

	if equipment.Name == "" {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", errors.New("error: empty equipment name"))
	}

	var dedodedJsonString map[string]string
	err = json.Unmarshal(c.Body(), &dedodedJsonString)
	if err != nil {
		log.Println("error decoding json to string")
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error decoding json to string", errors.New("error: error decoding json to string"))
	}
	// convert json strings to struc field types before create
	for key, value := range dedodedJsonString {
		if key == "servicing_period" {
			equipment.ServicingPeriod, err = strconv.Atoi(value)
			if err != nil {
				log.Println("error converting servicing period string to int")
				return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
					"eerror converting servicing period string to int", errors.New("error: error converting servicing period string to int"))
			}
		}
		if key == "hospital_id" {
			equipment.HospitalID, err = uuid.ParseBytes([]byte(value))
			if err != nil {
				log.Println("error converting hospital_id string to uuid")
				return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
					"error converting hospital_id string to uuid", errors.New("error: error converting hospital_id string to uuid"))
			}
		}
	}

	err = tx.Create(&equipment).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error commit creating equipment", err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	return apis.GeneralApiResponse(c, apis.StatusCreatedResponseCode, "successfully added equipment", &equipment)
}

func (h *Handler) UpdateEquipment(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	id := c.Params("id")

	equipment := &models.Equipment{}
	err := json.Unmarshal(c.Body(), equipment)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", err.Error())
	}

	if equipment.Name == "" {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", errors.New("error: empty equipment name"))
	}
	err = tx.Where("id = ?", id).Updates(equipment).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction on equipment update", err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	log.Println("Updating Equipment ...")

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully uptdated equipment", &equipment)
}

func (h *Handler) DeleteEquipment(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	id := c.Params("id")

	equipment := &models.Equipment{}
	err := tx.Where("id = ?", id).Delete(equipment).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction on equipment update", err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully deleted equipment", nil)
}
