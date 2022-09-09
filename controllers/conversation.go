package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

	query_max := db.Table("messages").Select("MAX(messages.timestamp) as latest").Group("conversation_id")
	query_count := db.Table("messages").Select("conversation_id,COUNT(`read`) as unread").Group("conversation_id").Where("sender_id != ? AND `read` = ?", ctx.User.ID, "0")
	if err := db.
		Select(`conversations.*,IF(msg.unread = NULL,0,msg.unread) as unread`).
		Preload("Users").
		Preload("Messages", func(g *gorm.DB) *gorm.DB {
			return g.
				Preload("Sender").
				Joins("JOIN (?) m ON m.latest = messages.timestamp", query_max)
		}).
		Joins("JOIN conversations_users on conversations_users.conversation_id = conversations.id").
		Joins("JOIN conversations_users c on c.conversation_id = conversations.id").
		Joins("JOIN users on conversations_users.users_id = users.id").
		Joins("JOIN messages on messages.conversation_id = conversations.id").
		Joins("LEFT JOIN (?) msg ON msg.conversation_id = messages.conversation_id", query_count).
		Where("conversations_users.users_id = ?", ctx.User.ID).
		Or("c.users_id = ?", ctx.User.ID).
		Group("conversations.id").
		Limit(int(query.PageSize)).
		Offset(int(query.Start)).
		Find(&conv).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Conversation not found"})
	}

	db.
		Table("conversations").
		Select(`count(*)`).
		Joins("JOIN conversations_users on conversations_users.conversation_id = conversations.id").
		Joins("JOIN conversations_users c on c.conversation_id = conversations.id").
		Joins("JOIN users on conversations_users.users_id = users.id").
		Where("conversations_users.users_id = ?", ctx.User.ID).
		Or("c.users_id = ?", ctx.User.ID).
		Group("conversations.id").
		Count(&query.Total)

	query.TotalPage = int64(math.Ceil(float64(query.Total) / float64(query.PageSize)))

	data := make([]models.ConversationWithLastMessage, len(conv))
	for i := range conv {
		var last models.Messages
		if len(conv[i].Messages) > 0 {
			last = conv[i].Messages[0]
		}
		data[i] = models.ConversationWithLastMessage{
			Conversation: conv[i],
			LastMessage:  last,
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}
