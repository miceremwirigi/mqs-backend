package models

import (
	"github.com/google/uuid"
)

type Equipment struct {
	BaseModel
	Name            string
	Model           string
	ServicingPeriod int
	HospitalID      uuid.UUID
	Hospital        Hospital   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Services        []*Service `gorm:"many2many:serviced_equipments;"`
}
