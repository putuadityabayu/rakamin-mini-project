package models

type Convertaion struct {
	ID    uint64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Users []Users `json:"users" gorm:"many2many:users"`
}
