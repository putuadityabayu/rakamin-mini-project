package models

import "rakamin.com/project/config"

func SetupModels() {
	config.DB.AutoMigrate(&Users{}, &Conversation{}, &Messages{})
}
