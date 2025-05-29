package models

import "time"

type Service struct {
	BaseModel
	Date       time.Time
	Equipments []*Equipment `gorm:"many2many:serviced_equipments;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
	Engineers  []*Engineer  `gorm:"many2many:serviced_by;constraint:OnUpdate:CASCADE;onDelete:SET NULL;"`
}
