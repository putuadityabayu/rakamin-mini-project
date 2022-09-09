package models

type Conversation struct {
	ID       uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Unread   uint64     `json:"unread" gorm:"<-:false;-:migration"`
	Users    []Users    `json:"users" gorm:"many2many:conversations_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Messages []Messages `json:"-" gorm:"foreignKey:ConversationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ConversationWithLastMessage struct {
	Conversation
	LastMessage Messages `json:"message" gorm:"-"`
}
