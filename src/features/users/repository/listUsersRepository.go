package repository

import (
	"context"
	_interface "main/features/users/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewListUsersRepository(gormDB *gorm.DB) _interface.IListUsersRepository {
	return &ListUsersRepository{GormDB: gormDB}
}
func (d *ListUsersRepository) FindUsers(ctx context.Context) ([]mysql.Users, error) {
	users := make([]mysql.Users, 0)
	err := d.GormDB.WithContext(ctx).Model(&users).Find(&users).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMongoDB)
	}
	return users, nil
}

func (d *ListUsersRepository) CountUsers(ctx context.Context) (int, error) {
	var count int64
	err := d.GormDB.WithContext(ctx).Model(&mysql.Users{}).Count(&count).Error
	if err != nil {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMongoDB)
	}
	return int(count), nil
}
