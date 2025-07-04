package models

import (
	"time"

	"github.com/google/uuid"
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
