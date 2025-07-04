package equipments

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/miceremwirigi/mqs-backend/common"
	"github.com/miceremwirigi/mqs-backend/common/apis"
	"github.com/miceremwirigi/mqs-backend/models"
)

func (h *Handler) GetAllEquipments(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": tx.Error.Error()})
	}
	var equipments []models.Equipment
	err := tx.Preload("Hospital").Preload("Department").Preload("Services").Find(&equipments).Error
	if err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": equipments})
}

func (h *Handler) GetEquipment(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var equipment models.Equipment
	err := tx.Preload("Hospital").Preload("Department").Preload("Services").Where("id = ?", id).Find(&equipment).Error
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
	err := tx.Preload("Hospital").Preload("Department").Preload("Services").Where("id = ?", id).Find(&equipment).Error
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
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error decoding json to string", errors.New("error: error decoding json to string"))
	}
	// convert json strings to struct field types before create
	for key, value := range dedodedJsonString {
		if key == "servicing_period" {
			equipment.ServicingPeriod, err = strconv.Atoi(value)
			if err != nil {
				return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
					"eerror converting servicing period string to int", errors.New("error: error converting servicing period string to int"))
			}
		}
		if key == "hospital_id" {
			equipment.HospitalID, err = uuid.ParseBytes([]byte(value))
			if err != nil {
				return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
					"error converting hospital_id string to uuid", errors.New("error: error converting hospital_id string to uuid"))
			}
		}
		if key == "department_id" {
			equipment.DepartmentID, err = uuid.ParseBytes([]byte(value))
			if err != nil {
				return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
					"error converting department_id string to uuid", errors.New("error: error converting department_id string to uuid"))
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

	// Fetch the existing equipment to update associations
	var existingEquipment models.Equipment
	if err := tx.Preload("Services").Preload("Department").First(&existingEquipment, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"equipment not found", err.Error())
	}

	equipment := &models.Equipment{}
	err := json.Unmarshal(c.Body(), equipment)
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", err.Error())
	}

	if equipment.Name == "" {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", errors.New("error: empty equipment name"))
	}

	// Update scalar fields
	if err := tx.Where("id = ?", id).Updates(equipment).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction on equipment update", err.Error())
	}

	// If services are provided in the update, update the association
	var dedodedJsonString map[string]string
	err = json.Unmarshal(c.Body(), &dedodedJsonString)
	if err == nil {
		if value, ok := dedodedJsonString["services"]; ok {
			serviceIDs := strings.Split(value, ",")
			services := []*models.Service{}
			for _, service_id := range serviceIDs {
				service := models.Service{}
				tx.First(&service, "id = ?", service_id)
				services = append(services, &service)
			}
			if err := tx.Model(&existingEquipment).Association("Services").Replace(services); err != nil {
				tx.Rollback()
				return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
					"error updating services association", err.Error())
			}
		}
	}
	// Update department if provided
	if value, ok := dedodedJsonString["department_id"]; ok {
		equipment.DepartmentID, err = uuid.ParseBytes([]byte(value))
		if err != nil {
			tx.Rollback()
			return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
				"error converting department_id string to uuid", errors.New("error: error converting department_id string to uuid"))
		}
		if err := tx.Where("id = ?", id).Updates(equipment).Error; err != nil {
			tx.Rollback()
			return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
				"error updating department association", err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully uptdated equipment", equipment)
}

func (h *Handler) DeleteEquipment(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	id := c.Params("id")

	// Remove associations in serviced_equipments before deleting equipment
	if err := tx.Exec("DELETE FROM serviced_equipments WHERE equipment_id = ?", id).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error deleting related serviced_equipments", err.Error())
	}

	equipment := &models.Equipment{}
	err := tx.Where("id = ?", id).Delete(equipment).Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction on equipment update", err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully deleted equipment", nil)
}

func (h *Handler) SetEquipmentDone(c fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}
	return h.DB.Model(&models.Equipment{}).Where("id = ?", id).Update("is_done", true).Error
}

func (h *Handler) SetEquipmentDoNotRemind(c fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}
	return h.DB.Model(&models.Equipment{}).Where("id = ?", id).Update("snooze_email", true).Error
}

func (h *Handler) SnoozeEquipmentReminder(c fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing id"})
	}
	now := time.Now()
	return h.DB.Model(&models.Equipment{}).Where("id = ?", id).Update("last_reminder_date", now).Error
}
