package models

type Conversation struct {
	ID       uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Users    []Users    `json:"users" gorm:"many2many:conversations_users"`
	Messages []Messages `json:"-" gorm:"foreignKey:ConversationID;references:ID"`
}

type ConversationWithLastMessage struct {
	Conversation
	LastMessage Messages `json:"last_message" gorm:"-"`
}
