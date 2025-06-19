package models

type Department struct {
	BaseModel
	Name        string `json:"name" gorm:"not null;unique"`
	Description string
}
