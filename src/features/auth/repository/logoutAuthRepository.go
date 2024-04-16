package repository

import (
	"context"
	_interface "main/features/auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewLogoutAuthRepository(gormDB *gorm.DB) _interface.ILogoutAuthRepository {
	return &LogoutAuthRepository{GormDB: gormDB}
}

func (d *LogoutAuthRepository) FindOneAndUpdateUser(ctx context.Context, uID uint) error {
	user := mysql.Users{
		State:  "logout",
		RoomID: 1,
	}
	result := d.GormDB.WithContext(ctx).Model(&user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)
	}
	return nil
}
