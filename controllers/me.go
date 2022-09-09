package controllers

import (
	"github.com/gofiber/fiber/v2"
	"rakamin.com/project/models"
)

func Me(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(*models.Context)

	return c.Status(fiber.StatusOK).JSON(ctx.User)
}
