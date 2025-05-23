package models

import "time"

type Service struct {
	BaseModel
	Date       time.Time
	Equipments []Equipment `gorm:"many2many:serviced_equipments;"`
	Engineers  []Engineer  `gorm:"many2many:serviced_by;"`
}
