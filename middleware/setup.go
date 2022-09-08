package middleware

import (
	"github.com/gofiber/fiber/v2"
	"rakamin.com/project/models"
)

// Set context default value
func Setup(c *fiber.Ctx) error {
	c.Locals("ctx", &models.Context{})
	return c.Next()
}
