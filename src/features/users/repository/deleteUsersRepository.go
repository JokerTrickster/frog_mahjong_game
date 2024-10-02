package repository

import (
	"context"
	_interface "main/features/users/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewDeleteUsersRepository(gormDB *gorm.DB) _interface.IDeleteUsersRepository {
	return &DeleteUsersRepository{GormDB: gormDB}
}

func (d *DeleteUsersRepository) FindOneAndDeleteUsers(ctx context.Context, userID uint) error {
	err := d.GormDB.WithContext(ctx).Where("id = ?", userID).Delete(&mysql.Users{}).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), userID), utils.ErrFromMongoDB)
	}
	return nil
}
