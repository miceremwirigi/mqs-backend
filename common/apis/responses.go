package apis

import "github.com/gofiber/fiber/v3"

var (
	StatusOkResponseCode                  = 200
	StatusBadRequestResponseCode          = 400
	StatusInternalServerErrorResponseCode = 500
	StatusForbiddenResponseCode           = 403
	StatusNotFoundResponseCode            = 404
	StatusCreatedResponseCode             = 201
	StatusNoContentResponseCode           = 204
	StatusUnauthorizedResponseCode        = 401
	StatusConflictResponseCode            = 409 // Added for duplicate/unique constraint
)

func GeneralApiResponse(c fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode,
		"message": message,
		"data":    data,
	})
}
