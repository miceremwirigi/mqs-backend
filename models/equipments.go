package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Equipment struct {
	BaseModel
	Name             string
	Model            string
	ServicingPeriod  int
	HospitalID       uuid.UUID
	Hospital         Hospital `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DepartmentID     uuid.UUID
	Department       Department `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SerialNumber     string
	Manufacturer     string
	Services         []*Service `gorm:"many2many:serviced_equipments;"`
	LastReminderDate *time.Time `gorm:"column:last_reminder_date"`
	SnoozeEmail      bool       `gorm:"column:snooze_email;default:false"`
}

// IsDone checks if the equipment service is marked as done
func (eq Equipment) IsDone() bool {
	for _, service := range eq.Services {
		if service.Date.AddDate(0, eq.ServicingPeriod, 0).After(time.Now()) {
			return true
		}
	}
	return false
}

// DueDate calculates the next due date for servicing the equipment
// based on the last service date and servicing period.
// If past, it returns the date when the equipment was due for service
func (eq Equipment) DueDate() time.Time {
	overdue := eq.CreatedAt
	for _, service := range eq.Services {
		if service.Date.AddDate(0, eq.ServicingPeriod, 0).After(time.Now()) {
			return service.Date.AddDate(0, eq.ServicingPeriod, 0)
		} else if service.Date.AddDate(0, eq.ServicingPeriod, 0).After(overdue) {
			overdue = service.Date.AddDate(0, eq.ServicingPeriod, 0)
		}
	}
	return overdue
}

// Retreives the service engiuneers email for the equipment
func (eq Equipment) EngineersEmail() string {
	for _, service := range eq.Services {
		if service.Engineers != nil {
			for _, engineer := range service.Engineers {
				if engineer.Email != "" {
					return engineer.Email
				}
			}
		}
	}
	return ""
}

func GetEquipments(db *gorm.DB) ([]Equipment, error) {
	var equipments []Equipment
	if err := db.Preload("Hospital").Preload("Department").Preload("Services.Engineers").Find(&equipments).Error; err != nil {
		return nil, err
	}
	return equipments, nil
}

func GetEquipmentByID(db *gorm.DB, id string) (*Equipment, error) {
	var equipment Equipment
	if err := db.Preload("Hospital").Preload("Department").Preload("Services.Engineers").First(&equipment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &equipment, nil
}
