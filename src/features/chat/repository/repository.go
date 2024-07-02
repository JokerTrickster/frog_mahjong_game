package repository

import "gorm.io/gorm"

type MessageChatRepository struct {
	GormDB *gorm.DB
}

type AuthChatRepository struct {
	GormDB *gorm.DB
}
