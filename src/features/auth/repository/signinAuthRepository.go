package repository

import (
	"context"
	_errors "main/features/auth/model/errors"
	_interface "main/features/auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewSigninAuthRepository(gormDB *gorm.DB) _interface.ISigninAuthRepository {
	return &SigninAuthRepository{GormDB: gormDB}
}
func (g *SigninAuthRepository) FindOneUser(ctx context.Context, email, password string) (mysql.Users, error) {
	var user mysql.Users
	result := g.GormDB.WithContext(ctx).Where("email = ? AND password = ?", email, password).First(&user)
	if result.RowsAffected == 0 {
		return mysql.Users{}, utils.ErrorMsg(ctx, utils.ErrUserNotExist, utils.Trace(), _errors.ErrUserNotFound.Error(), utils.ErrFromClient)
	}
	if result.Error != nil {
		return mysql.Users{}, result.Error
	}
	return user, nil
}
