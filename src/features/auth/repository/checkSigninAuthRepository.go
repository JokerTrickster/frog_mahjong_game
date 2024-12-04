package repository

import (
	"context"
	_errors "main/features/auth/model/errors"
	_interface "main/features/auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewCheckSigninAuthRepository(gormDB *gorm.DB) _interface.ICheckSigninAuthRepository {
	return &CheckSigninAuthRepository{GormDB: gormDB}
}

func (g *CheckSigninAuthRepository) FindOneAndUpdateUser(ctx context.Context, email, password string) (mysql.Users, error) {
	user := mysql.Users{
		Email:  email,
		State:  "wait",
		RoomID: 1,
	}
	// 사용자 정보를 가져온다.
	var findUser mysql.Users
	err := g.GormDB.WithContext(ctx).Model(&findUser).Where("email = ? ", email).First(&findUser).Error
	if err != nil {
		return mysql.Users{}, utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), utils.HandleError(_errors.ErrUserNotFound.Error()+err.Error(), email, password), utils.ErrFromClient)
	}
	if password != findUser.Password {
		return mysql.Users{}, utils.ErrorMsg(ctx, utils.ErrPasswordNotMatch, utils.Trace(), utils.HandleError(_errors.ErrPasswordNotMatch.Error(), password), utils.ErrFromClient)
	}

	result := g.GormDB.WithContext(ctx).Model(&user).Where("email = ?", email).Updates(user)
	if result.Error != nil {
		return mysql.Users{}, utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), utils.HandleError(_errors.ErrUserNotFound.Error()+result.Error.Error(), email, password), utils.ErrFromClient)
	}

	// 변경된 사용자 정보를 가져옵니다.
	err = g.GormDB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return mysql.Users{}, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(err.Error(), email, password), utils.ErrFromInternal)
	}
	return user, nil
}

func (d *CheckSigninAuthRepository) CheckToken(ctx context.Context, uID uint) (*mysql.Tokens, error) {
	token := mysql.Tokens{
		UserID: uID,
	}
	err := d.GormDB.WithContext(ctx).Model(&token).Where("user_id = ?", uID).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(err.Error(), uID), utils.ErrFromInternal)
	}
	return &token, nil
}
