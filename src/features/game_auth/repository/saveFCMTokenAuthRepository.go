package repository

import (
	"context"
	_interface "main/features/game_auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewFCMTokenAuthRepository(gormDB *gorm.DB) _interface.IFCMTokenAuthRepository {
	return &FCMTokenAuthRepository{GormDB: gormDB}
}

func (d *FCMTokenAuthRepository) SaveFCMToken(ctx context.Context, userID uint, token string) error {
	// Define the user token structure
	userToken := mysql.UserTokens{
		UserID: userID,
		Token:  token,
	}

	// Use Gorm's upsert feature to insert or update
	err := d.GormDB.WithContext(ctx).
		Model(&mysql.UserTokens{}).
		Where("user_id = ?", userID).
		Assign(userToken).
		FirstOrCreate(&userToken).Error

	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMongoDB)
	}

	return nil
}
