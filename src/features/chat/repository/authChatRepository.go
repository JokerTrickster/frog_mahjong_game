package repository

import (
	"context"
	_interface "main/features/chat/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewAuthChatRepository(gormDB *gorm.DB) _interface.IAuthChatRepository {
	return &AuthChatRepository{GormDB: gormDB}
}

func (d *AuthChatRepository) InsertOneChat(ctx context.Context, chatDTO *mysql.Chats) error {
	err := d.GormDB.Create(&chatDTO).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(),chatDTO), utils.ErrFromMysqlDB)
	}
	return nil
}

func (d *AuthChatRepository) FindOneUserInfo(ctx context.Context, userID uint) (*mysql.Users, error) {
	var user *mysql.Users
	err := d.GormDB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(),userID), utils.ErrFromMysqlDB)
	}
	return user, nil
}
