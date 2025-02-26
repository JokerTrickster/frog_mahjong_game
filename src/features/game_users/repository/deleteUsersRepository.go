package repository

import (
	"context"
	_interface "main/features/game_users/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewDeleteUsersRepository(gormDB *gorm.DB) _interface.IDeleteUsersRepository {
	return &DeleteUsersRepository{GormDB: gormDB}
}

func (d *DeleteUsersRepository) FindOneAndDeleteUsers(ctx context.Context, userID uint) error {
	err := d.GormDB.WithContext(ctx).Where("id = ?", userID).Delete(&mysql.GameUsers{}).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), userID), utils.ErrFromMongoDB)
	}
	return nil
}

func (d *DeleteUsersRepository) DeleteToken(ctx context.Context, userID uint) error {
	token := mysql.Tokens{
		UserID: userID,
	}
	result := d.GormDB.Model(&token).Where("user_id = ?", userID).Delete(&token)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(result.Error.Error(), userID), utils.ErrFromInternal)
	}
	return nil
}
