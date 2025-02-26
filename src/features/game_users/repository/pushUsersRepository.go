package repository

import (
	"context"
	_interface "main/features/game_users/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewPushUsersRepository(gormDB *gorm.DB) _interface.IPushUsersRepository {
	return &PushUsersRepository{GormDB: gormDB}
}

func (d *PushUsersRepository) FindUsersForNotifications(c context.Context) ([]mysql.GameUsers, error) {
	var users []mysql.GameUsers
	err := d.GormDB.Find(&users).Where("alert_enabled = ?", true).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (d *PushUsersRepository) FindOnePushToken(c context.Context, userID uint) (string, error) {
	var userTokens mysql.UserTokens
	err := d.GormDB.Find(&userTokens).Where("user_id = ?", userID).Error
	if err != nil {
		return "", err
	}
	return userTokens.Token, nil
}
