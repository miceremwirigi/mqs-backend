package auth

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/miceremwirigi/mqs-backend/models"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type loginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateHashFromPassword(password string) (string, error) {
	bcryptCost := bcrypt.DefaultCost
	if bcryptCost < 4 || bcryptCost > 31 {
		return "", errors.New("bcrypt cost out of range: must be between 4 and 31")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	hashedPasswordString := string(hashedPassword)
	return hashedPasswordString, nil // Replace with actual hashing logic
}

func (h *Handler) LoginUser(c fiber.Ctx) (string, error) {
	userInfo := new(loginInfo)
	if err := json.Unmarshal(c.Body(), userInfo); err != nil {
		return "", errors.New("failed to parse request body")
	}
	if userInfo == nil {
		return "", errors.New("request body cannot be empty")
	}
	username := userInfo.Username
	password := userInfo.Password
	if username == "" || password == "" {
		return "", errors.New("username and password cannot be empty")
	}

	var user models.User
	tx := h.DB.Begin()
	if tx.Error != nil {
		return "", errors.New("failed to begin transaction")
	}
	err := tx.Where("username = ?", username).First(&user).Error
	if err != nil {
		tx.Rollback()
		return "", errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		tx.Rollback()
		return "", errors.New("invalid password"+": " + user.PasswordHash)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return "", errors.New("error committing transaction")
	}

	var jwtKey = []byte("your_secret_key") // Replace with your actual secret key

	claims := Claims{
		Username: user.Username,
		UserID:   (user.ID).String(),
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenObj.SignedString(jwtKey)
	if err != nil {
		tx.Rollback()
		return "", errors.New("failed to generate token")
	}
	return token, nil
}

// Client side token deletion from storage already implemented
func (h *Handler) LogoutUser(c fiber.Ctx) error {
	// Invalidate the token by removing it from the client side or using a blacklist
	return c.SendStatus(fiber.StatusOK)
}
