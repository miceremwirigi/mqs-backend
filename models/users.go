package models

// User represents a user in the system
type User struct {
	BaseModel
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Email        string `gorm:"unique;not null"`
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Role         string `gorm:"not null"` // e.g., "admin", "engineer", "user"
}
