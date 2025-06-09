package models

// User represents a user in the system
type User struct {
	BaseModel
	Username     string `gorm:"unique;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"password"`
	Email        string `gorm:"unique;not null" json:"email"`
	FirstName    string `gorm:"not null" json:"first_name"`
	LastName     string `gorm:"not null" json:"last_name"`
	Role         string `gorm:"not null" json:"role"` // e.g., "admin", "engineer", "user"
}
