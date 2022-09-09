package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
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

	var conv models.Conversation
	var msg models.Messages
	var databaseErr error

	// Check conversation exists or created
	conv_where := models.Conversation{
		Users: []models.Users{
			ctx.User,
			user,
		},
	}

	if err := db.Preload("Users").
		Select("conversations.*,'0' as unread").
		Joins("JOIN conversations_users on conversations_users.conversation_id = conversations.id").
		Joins("JOIN conversations_users c on c.conversation_id = conversations.id").
		Joins("JOIN users on conversations_users.users_id = users.id").
		Where("conversations_users.users_id = ? AND c.users_id = ?", ctx.User.ID, user.ID).
		Or("conversations_users.users_id = ? AND c.users_id = ?", user.ID, ctx.User.ID).
		First(&conv).Error; err != nil {
		conv = conv_where
		msg = models.Messages{
			ConversationInternal: conv,
			Messages:             post.Message,
			Sender:               ctx.User,
		}
		msg.ConversationInternal.Unread = 1
		databaseErr = db.Save(&msg).Error
	} else {
		msg = models.Messages{
			ConversationID:       conv.ID,
			ConversationInternal: conv,
			Messages:             post.Message,
			SenderID:             ctx.User.ID,
			Sender:               ctx.User,
		}
		msg.ConversationInternal.Unread = 1
		databaseErr = db.Omit(clause.Associations).Create(&msg).Error
	}

	if databaseErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	result := models.MessagesWithConversation{
		Messages:     msg,
		Conversation: msg.ConversationInternal,
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func ReplyMessage(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(*models.Context)
	db := config.DB
	id_str := c.Params("id", "notfound")
	id_int, err := strconv.Atoi(id_str)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Conversation not found"})
	}
	id := uint64(id_int)

	var post models.PostNewMessage
	c.BodyParser(&post)

	// Check Message
	if post.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Message cannot be empty"})
	}

	// Check conversations
	var conv models.Conversation
	if err := db.Preload("Users").
		Select("conversations.*,'0' as unread").
		Joins("JOIN conversations_users on conversations_users.conversation_id = conversations.id").
		Joins("JOIN conversations_users c on c.conversation_id = conversations.id").
		Joins("JOIN users on conversations_users.users_id = users.id").
		Where("conversations.id = ?", id).
		First(&conv).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Conversation not found"})
	}

	msg := models.Messages{
		ConversationID:       conv.ID,
		ConversationInternal: conv,
		Messages:             post.Message,
		SenderID:             ctx.User.ID,
		Sender:               ctx.User,
	}
	msg.ConversationInternal.Unread = 1

	if err := db.Omit(clause.Associations).Create(&msg).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	result := models.MessagesWithConversation{
		Messages:     msg,
		Conversation: msg.ConversationInternal,
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func ListMessages(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(*models.Context)
	db := config.DB
	id_str := c.Params("id", "notfound")
	id_int, err := strconv.Atoi(id_str)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Conversation not found"})
	}
	id := uint64(id_int)

	var query models.Pagination
	c.QueryParser(&query)
	query.Format()

	var msg []models.Messages

	db.Model(&models.Messages{}).Limit(int(query.PageSize)).Offset(int(query.Start)).Where("conversation_id = ? AND sender_id != ?", id, ctx.User.ID).Update("read", true)

	db.
		Select(`messages.id,messages.conversation_id,messages.timestamp,messages.sender_id,messages.messages,IF(messages.sender_id = ?,'1',messages.read) as 'read'`, ctx.User.ID).
		Joins("JOIN conversations c on c.id = messages.conversation_id").
		Joins("JOIN conversations_users cu on cu.conversation_id = c.id").
		Joins("JOIN conversations_users cuu on cuu.conversation_id = c.id").
		Where("c.id = ? AND cu.users_id = ?", id, ctx.User.ID).
		Or("c.id = ? AND cuu.users_id = ?", id, ctx.User.ID).
		Group("messages.id").
		Order("messages.timestamp DESC").
		Limit(int(query.PageSize)).Offset(int(query.Start)).Preload("Sender").Find(&msg)

	db.Table("messages").Select("count(*)").
		Joins("JOIN conversations c on c.id = messages.conversation_id").
		Joins("JOIN conversations_users cu on cu.conversation_id = c.id").
		Joins("JOIN conversations_users cuu on cuu.conversation_id = c.id").
		Where("c.id = ? AND cu.users_id = ?", id, ctx.User.ID).
		Or("c.id = ? AND cuu.users_id = ?", id, ctx.User.ID).
		Group("messages.id").
		Count(&query.Total)
	query.TotalPage = int64(math.Ceil(float64(query.Total) / float64(query.PageSize)))

	result := models.PaginationResponse[models.Messages]{
		Pagination: query,
		Data:       msg,
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
