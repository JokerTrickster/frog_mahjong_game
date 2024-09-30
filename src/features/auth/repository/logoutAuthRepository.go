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
func (d *LogoutAuthRepository) DeleteToken(ctx context.Context, uID uint) error {
	token := mysql.Tokens{
		UserID: uID,
	}
	result := d.GormDB.Model(&token).Where("user_id = ?", uID).Delete(&token)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(result.Error.Error(),uID), utils.ErrFromInternal)
	}
	return nil
}

func (d *LogoutAuthRepository) FindOneAndUpdateUser(ctx context.Context, uID uint) error {
	user := mysql.Users{
		State:  "logout",
		RoomID: 1,
	}
	result := d.GormDB.WithContext(ctx).Model(&user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), utils.HandleError(result.Error.Error(),uID), utils.ErrFromClient)
	}
	return nil
}
