package repository

import (
	"context"
	_interface "main/features/chat/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewMessageChatRepository(gormDB *gorm.DB) _interface.IMessageChatRepository {
	return &MessageChatRepository{GormDB: gormDB}
}

func (d *MessageChatRepository) FindOneChat(ctx context.Context, secret string) (*mysql.Chat, error) {
	var chat mysql.Chat
	err := d.GormDB.WithContext(ctx).Where("secret = ?", secret).First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}
