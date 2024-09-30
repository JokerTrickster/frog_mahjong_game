package repository

import (
	"context"
	_interface "main/features/users/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewGetUsersRepository(gormDB *gorm.DB) _interface.IGetUsersRepository {
	return &GetUsersRepository{GormDB: gormDB}
}

func (d *GetUsersRepository) FindOneUser(ctx context.Context, userID int) (mysql.Users, error) {
	var userDTO mysql.Users
	err := d.GormDB.WithContext(ctx).Model(&userDTO).Where("id = ?", userID).First(&userDTO).Error
	if err != nil {
		return mysql.Users{}, utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(),utils.HandleError( err.Error()), utils.ErrFromClient)
	}
	return userDTO, nil
}
