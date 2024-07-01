package repository

import (
	"context"
	_interface "main/features/chat/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewAuthChatRepository(gormDB *gorm.DB) _interface.IAuthChatRepository {
	return &AuthChatRepository{GormDB: gormDB}
}

func (d *AuthChatRepository) InsertOneChat(ctx context.Context, chatDTO *mysql.Chat) error {
	err := d.GormDB.Create(&chatDTO).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *AuthChatRepository) FindOneUserInfo(ctx context.Context, userID uint) (*mysql.Users, error) {
	var user *mysql.Users
	err := d.GormDB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
