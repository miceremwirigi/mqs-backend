package micerejwt

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	SigningKey    []byte
	ContextKey    string // e.g. "user"
	SigningMethod string // e.g. "HS256"
}

func New(cfg Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if cfg.SigningMethod != "" && token.Method.Alg() != cfg.SigningMethod {
				return nil, fiber.ErrUnauthorized
			}
			return cfg.SigningKey, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Store token in context if ContextKey is set
		if cfg.ContextKey != "" {
			c.Locals(cfg.ContextKey, token)
		}
		return c.Next()
	}
}
