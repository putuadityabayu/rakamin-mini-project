package models

import "time"

type Messages struct {
	ID             uint64      `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID uint64      `json:"-"`
	Read           bool        `json:"read"`
	Timestamp      time.Time   `json:"timestamp" gorm:"autoCreateTime"`
	Conversation   Convertaion `json:"conversation"`
}
