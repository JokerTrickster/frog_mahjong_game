package repository

import (
	"context"
	_interface "main/features/game_users/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewListUsersRepository(gormDB *gorm.DB) _interface.IListUsersRepository {
	return &ListUsersRepository{GormDB: gormDB}
}
func (d *ListUsersRepository) FindUsers(ctx context.Context) ([]mysql.GameUsers, error) {
	users := make([]mysql.GameUsers, 0)
	err := d.GormDB.WithContext(ctx).Model(&users).Find(&users).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMongoDB)
	}
	return users, nil
}

func (d *ListUsersRepository) CountUsers(ctx context.Context) (int, error) {
	var count int64
	err := d.GormDB.WithContext(ctx).Model(&mysql.GameUsers{}).Count(&count).Error
	if err != nil {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMongoDB)
	}
	return int(count), nil
}
