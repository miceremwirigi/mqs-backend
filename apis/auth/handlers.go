package auth

import "gorm.io/gorm"

type Handler struct {
	DB *gorm.DB // Replace with actual DB type, e.g., *gorm.DB
}
