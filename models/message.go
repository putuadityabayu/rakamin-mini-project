package models

import "time"

type Messages struct {
	ID                   uint64       `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID       uint64       `json:"-" gorm:"not null"`
	SenderID             uint64       `json:"-" gorm:"not null"`
	Read                 bool         `json:"read" gorm:"default:false"`
	Timestamp            time.Time    `json:"timestamp" gorm:"autoCreateTime"`
	Messages             string       `json:"message" gorm:"not null"`
	ConversationInternal Conversation `json:"-" gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Sender               Users        `json:"sender" gorm:"foreignKey:SenderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type MessagesWithConversation struct {
	Messages
	Conversation Conversation `json:"conversation"`
}

type PostNewMessage struct {
	UserID  uint64 `json:"user_id"`
	Message string `json:"message"`
}
