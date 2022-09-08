package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rakamin.com/project/config"
	"rakamin.com/project/models"
)

func ListConversation(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(*models.Context)
	db := config.DB
	var query models.Pagination
	c.QueryParser(&query)
	query.Format()

	var conv []models.Conversation

	if err := db.Limit(int(query.PageSize)).Offset(int(query.Start)).Preload("Users").Preload("Messages", func(g *gorm.DB) *gorm.DB {
		return g.Order("timestamp DESC").Limit(1)
	}).Find(&conv, &models.Conversation{Users: []models.Users{ctx.User}}).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Conversation not found"})
	}
	db.Where(&models.Conversation{Users: []models.Users{ctx.User}}).Preload(clause.Associations).Count(&query.Total)
	query.TotalPage = int64(math.Ceil(float64(query.Total) / float64(query.PageSize)))

	data := make([]models.ConversationWithLastMessage, len(conv))
	for i := range conv {
		data[i] = models.ConversationWithLastMessage{
			Conversation: conv[i],
			LastMessage:  conv[i].Messages[0],
		}
	}

	result := models.PaginationResponse[models.ConversationWithLastMessage]{
		Pagination: query,
		Data:       data,
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func DeleteConversation(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(*models.Context)
	db := config.DB
	id_str := c.Params("id", "notfound")
	id_int, err := strconv.Atoi(id_str)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Conversation not found"})
	}

	id := uint64(id_int)

	if err := db.Select("Users", "Messages").Delete(&models.Conversation{ID: id, Users: []models.Users{ctx.User}}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": true})
}
