package controllers

import (
	"github.com/gofiber/fiber/v2"
	"rakamin.com/project/config"
	"rakamin.com/project/models"
)

func NewMessage(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(*models.Context)
	db := config.DB
	var post models.PostNewMessage
	var user models.Users
	c.BodyParser(&post)

	// Check POST
	if post.UserID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing user_id"})
	}

	// Check Message
	if post.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Message cannot be empty"})
	}

	// Check if post.user_id != my.id
	if ctx.User.ID == post.UserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id. Cannot chat yourself"})
	}

	// Check if post.user_id exist
	if err := db.First(&user, post.UserID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id. User not found"})
	}

	msg := models.Messages{
		ConversationInternal: models.Conversation{
			Users: []models.Users{
				ctx.User,
				user,
			},
		},
		Read:     false,
		Messages: post.Message,
	}

	db.Create(&msg)

	result := models.MessagesWithConversation{
		Messages:     msg,
		Conversation: msg.ConversationInternal,
	}

	return c.Status(fiber.StatusNotFound).JSON(result)
}
