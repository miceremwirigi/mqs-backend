package users

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/miceremwirigi/mqs-backend/apis/auth"
	"github.com/miceremwirigi/mqs-backend/common/apis"
	"github.com/miceremwirigi/mqs-backend/models"
)

func (h *Handler) GetUserByEmail(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	email := c.Params("email")

	var user models.User
	err := tx.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving user", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived user datail", user)
}

func (h *Handler) GetUserById(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var user models.User
	err := tx.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving user", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived user datail", user)
}

func (h *Handler) GetUserByUsername(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	username := c.Params("username")

	var user models.User
	err := tx.Where("username = ?", username).Find(&user).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving user", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived user datail", user)
}

func (h *Handler) GetAllUsers(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	var users []models.User
	err := tx.Find(&users).Error
	if err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving users", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "successfully retreived users", users)
}

func (h *Handler) AddUser(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	// Check if any users exist in the database
	var userCount int64
	if err := tx.Model(&models.User{}).Count(&userCount).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to count users", err)
	}

	var user models.User
	if err := json.Unmarshal(c.Body(), &user); err != nil {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error unmarshalling user data", err)
	}

	passwordHash, err := auth.GenerateHashFromPassword(user.PasswordHash)
	if err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error hashing password", err)
	}
	user.PasswordHash = passwordHash

	if user.Username == "" {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"username is required", nil)
	}
	if user.Email == "" {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"email is required", nil)
	}

	if userCount == 0 {
		// If this is the first user, force role to "admin"
		user.Role = "admin"
	} else {
		// For all other users, require the requester to be an admin
		// Only check for Authorization header if there is at least one user
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
				"authorization header missing", nil)
		}
		// Expect header in format "Bearer <token>"
		var tokenString string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
				"invalid authorization header format", nil)
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Replace with your JWT secret or key lookup
			return []byte(auth.JwtKey), nil
		})
		if err != nil || !token.Valid {
			return apis.GeneralApiResponse(c, apis.StatusForbiddenResponseCode,
				"invalid token", err)
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return apis.GeneralApiResponse(c, apis.StatusForbiddenResponseCode,
				"invalid token claims", nil)
		}
		if claims == nil {
			log.Println("JWT claims not found in context")
			return apis.GeneralApiResponse(c, apis.StatusForbiddenResponseCode,
				"only admins can create users", nil)
		}
		// Example for github.com/golang-jwt/jwt/v5
		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			log.Println("Unauthorized user creation attempt by non-admin")
			return apis.GeneralApiResponse(c, apis.StatusForbiddenResponseCode,
				"only admins can create users", nil)
		}
		// If no role provided, default to "user"
		if user.Role == "" {
			user.Role = "user"
		}
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error creating user", err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing transaction", err)
	}
	return apis.GeneralApiResponse(c, apis.StatusCreatedResponseCode, "user successfully registered", user)
}

func (h *Handler) UpdateUser(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	var user models.User
	var newUserInfo models.User
	if err := tx.Where("id = ?", id).First(&user).Error; err != nil {
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error retreiving user", err)
	}

	if err := json.Unmarshal(c.Body(), &newUserInfo); err != nil {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"error unmarshalling user data", err)
	}

	if newUserInfo.PasswordHash != "" || newUserInfo.Role != "" {
		return apis.GeneralApiResponse(c, apis.StatusBadRequestResponseCode,
			"sensitive data cannot be updated via this endpoint", nil)
	}

	if newUserInfo.Username != "" {
		user.Username = newUserInfo.Username
	}
	if newUserInfo.Email != "" {
		user.Email = newUserInfo.Email
	}
	if newUserInfo.FirstName != "" {
		user.FirstName = newUserInfo.FirstName
	}
	if newUserInfo.LastName != "" {
		user.LastName = newUserInfo.LastName
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error updating user", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "user successfully updated", user)
}

func (h *Handler) DeleteUser(c fiber.Ctx) error {
	tx := h.DB.Begin()
	if tx.Error != nil {
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"failed to begin transaction", nil)
	}

	id := c.Params("id")

	if err := tx.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusNotFoundResponseCode,
			"error deleting user", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apis.GeneralApiResponse(c, apis.StatusInternalServerErrorResponseCode,
			"error committing transaction", err)
	}

	return apis.GeneralApiResponse(c, apis.StatusOkResponseCode, "user successfully deleted", nil)
}
