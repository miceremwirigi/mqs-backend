package departments

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/miceremwirigi/mqs-backend/models"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

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
		return c.Status(500).JSON(fiber.Map{"error": tx.Error.Error()})
	}
	var dept models.Department
	if err := json.Unmarshal(c.Body(), &dept); err != nil {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dept.Name == "" {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Department name required"})
	}
	if err := tx.Create(&dept).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"data": dept})
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
