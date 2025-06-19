package departments

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/miceremwirigi/mqs-backend/common/apis"
	"github.com/miceremwirigi/mqs-backend/models"
	"gorm.io/gorm"
)

func (h *Handler) GetAllDepartments(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": tx.Error.Error()})
	}
	var departments []models.Department
	if err := tx.Find(&departments).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": departments})
}

func (h *Handler) GetDepartment(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": tx.Error.Error()})
	}
	id := c.Params("id")
	var department models.Department
	if err := tx.First(&department, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Department not found"})
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": department})
}

func (h *Handler) AddDepartment(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", tx.Error.Error())
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(c.Body(), &input); err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error binding body to struct", err.Error())
	}
	if input.Name == "" {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"department name required", "empty department name")
	}

	// Check for duplicate department name
	var existing models.Department
	if err := tx.Where("name = ?", input.Name).First(&existing).Error; err == nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusConflictResponseCode,
			"department name must be unique", "duplicate department name")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error checking for duplicate department", err.Error())
	}

	department := models.Department{
		Name: input.Name,
	}

	if err := tx.Create(&department).Error; err != nil {
		// Defensive: catch DB-level unique constraint error
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE") {
			tx.Rollback()
			return apis.GeneralApiResponse(c, apis.StatusConflictResponseCode,
				"department name must be unique", "duplicate department name")
		}
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error creating department", err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing transaction", err.Error())
	}

	return apis.GeneralApiResponse(c, apis.StatusCreatedResponseCode,
		"successfully added department", &department)
}

func (h *Handler) UpdateDepartment(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": tx.Error.Error()})
	}
	id := c.Params("id")
	var dept models.Department
	if err := tx.First(&dept, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Department not found"})
	}
	var update models.Department
	if err := json.Unmarshal(c.Body(), &update); err != nil {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if update.Name != "" {
		dept.Name = update.Name
	}
	if err := tx.Save(&dept).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": dept})
}

func (h *Handler) DeleteDepartment(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": tx.Error.Error()})
	}
	id := c.Params("id")
	if err := tx.Delete(&models.Department{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

func (h *Handler) GetDepartmentHtml(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.Status(500).SendString(tx.Error.Error())
	}
	id := c.Params("id")
	var department models.Department
	if err := tx.First(&department, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return c.Status(404).SendString("Department not found")
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Render("department", fiber.Map{"Department": department})
}
