package models

type Engineer struct {
	BaseModel
	Name     string     `json:"name"`
	Contact  string     `json:"contact"`
	Services []*Service `gorm:"many2many:serviced_by;"`
}
