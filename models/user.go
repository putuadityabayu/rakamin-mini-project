package models

import (
	"crypto/md5"
	"encoding/hex"

	"gorm.io/gorm"
)

type Users struct {
	ID       uint64 `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	Name     string `json:"name" gorm:"not null"`
	UserName string `json:"username" gorm:"column:username;not null"`
	Password string `json:"-" gorm:"not null"`
}

type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	hash := md5.Sum([]byte(u.Password))
	pass_str := hex.EncodeToString(hash[:])
	u.Password = pass_str
	return
}

func (u *Users) BeforeUpdate(tx *gorm.DB) (err error) {
	hash := md5.Sum([]byte(u.Password))
	pass_str := hex.EncodeToString(hash[:])
	u.Password = pass_str
	return
}
