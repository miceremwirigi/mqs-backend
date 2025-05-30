package services

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/miceremwirigi/mqs-backend/common"
	"github.com/miceremwirigi/mqs-backend/common/apis"
	"github.com/miceremwirigi/mqs-backend/models"
)

func (h *Handler) GetAllServices(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	var services []models.Service
	err := tx.Model(&services).Preload("Equipments").
		Preload("Equipments").Preload("Equipments.Hospital").Preload("Engineers").Find(&services).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error getting services", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived services", services)
}

func (h *Handler) GetService(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var service models.Service
	err := tx.Preload("Equipments").Preload("Equipments.Hospital").Where("id = ?", id).Find(&service).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving service", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived service datail", service)
}

func (h *Handler) GetServiceHtml(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var service models.Service
	err := tx.Preload("Equipments").Preload("Equipments.Hospital").
		Preload("Engineers").Where("id = ?", id).Find(&service).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving service", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	serviceMap, err := common.StructToMap(service)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error converting service struct to map", err)
	}

	responseFiberMap := fiber.Map(serviceMap)

	return c.Render("service", responseFiberMap)
}

func (h *Handler) AddService(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	service := &models.Service{}
	// equipments := []models.Equipment{}
	// endineers := []models.Engineers{}

	// err := json.Unmarshal(c.Body(), service)
	// if err != nil {
	// 	return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
	// 		"error binding body to struct", err.Error())
	// }

	var dedodedJsonString map[string]string
	err := json.Unmarshal(c.Body(), &dedodedJsonString)
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error decoding json to string", errors.New("error: error decoding json to string"))
	}
	// convert json strings to struc field types before create
	for key, value := range dedodedJsonString {
		switch key {
		case "date":
			service.Date, err = time.Parse("2006-01-02T15:04:05.999999999Z07:00", value)
			if err != nil {
				return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
					"error converting string to date", errors.New("error: cannot convert string to date"))
			}
		case "equipments":
			equipments := []*models.Equipment{}
			equipmentIDs := strings.Split(value, ",")
			for _, equipment_id := range equipmentIDs {
				equipment := models.Equipment{}
				tx.First(&equipment, "id = ?", equipment_id)
				equipments = append(equipments, &equipment)
			}
			service.Equipments = equipments
		case "engineers":
			engineers := []*models.Engineer{}
			engineerIDs := strings.Split(value, ",")
			for _, engineer_id := range engineerIDs {
				engineer := models.Engineer{}
				tx.First(&engineer, "id = ?", engineer_id)
				engineers = append(engineers, &engineer)
			}
			service.Engineers = engineers
		default:
			log.Println(key + ": " + value)
		}
	}

	if service.Date.IsZero() {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", errors.New("error: empty service date"))
	}

	err = tx.Create(service).Error
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

	return apis.GeneralApiResponse(c, apis.StatusCreatedResponseCode, "successfully added service", &service)
}

func (h *Handler) UpdateService(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	id := c.Params("id")

	// 1. Fetch the existing service
	var service models.Service
	if err := tx.Preload("Equipments").Preload("Engineers").First(&service, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"service not found", err.Error())
	}

	var dedodedJsonString map[string]string
	err := json.Unmarshal(c.Body(), &dedodedJsonString)
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error decoding json to string", errors.New("error: error decoding json to string"))
	}

	for key, value := range dedodedJsonString {
		switch key {
		case "date":
			service.Date, err = time.Parse("2006-01-02T15:04:05.999999999Z07:00", value)
			if err != nil {
				tx.Rollback()
				return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
					"error converting string to date", errors.New("error: cannot convert string to date"))
			}
		case "equipments":
			equipments := []*models.Equipment{}
			equipmentIDs := strings.Split(value, ",")
			for _, equipment_id := range equipmentIDs {
				equipment := models.Equipment{}
				tx.First(&equipment, "id = ?", equipment_id)
				equipments = append(equipments, &equipment)
			}
			// Use GORM's Association().Replace()
			if err := tx.Model(&service).Association("Equipments").Replace(equipments); err != nil {
				tx.Rollback()
				return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
					"error updating equipments", err.Error())
			}
		case "engineers":
			engineers := []*models.Engineer{}
			engineerIDs := strings.Split(value, ",")
			for _, engineer_id := range engineerIDs {
				engineer := models.Engineer{}
				tx.First(&engineer, "id = ?", engineer_id)
				engineers = append(engineers, &engineer)
			}
			// Use GORM's Association().Replace()
			if err := tx.Model(&service).Association("Engineers").Replace(engineers); err != nil {
				tx.Rollback()
				return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
					"error updating engineers", err.Error())
			}
		default:
			log.Println(key + ": " + value)
		}
	}

	if service.Date.IsZero() {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", errors.New("error: empty service name"))
	}

	// Save the updated date (and any other scalar fields)
	if err := tx.Model(&service).Updates(service).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction on service update", err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully updated service", &service)
}

func (h *Handler) DeleteService(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	id := c.Params("id")

	service := &models.Service{}
	// First, delete related records in the join tables to avoid foreign key constraint errors
	if err := tx.Exec("DELETE FROM serviced_equipments WHERE service_id = ?", id).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error deleting related serviced_equipments", err.Error())
	}
	if err := tx.Exec("DELETE FROM serviced_by WHERE service_id = ?", id).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error deleting related serviced_by", err.Error())
	}
	err := tx.Where("id = ?", id).Delete(service).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing database transaction on service update", err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err.Error())
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully deleted service", nil)
}
