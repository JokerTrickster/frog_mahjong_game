package repository

import (
	"context"
	"errors"
	_errors "main/features/auth/model/errors"
	_interface "main/features/auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewValidatePasswordAuthRepository(gormDB *gorm.DB) _interface.IValidatePasswordAuthRepository {
	return &ValidatePasswordAuthRepository{GormDB: gormDB}
}

func (g *ValidatePasswordAuthRepository) CheckAuthCode(ctx context.Context, email, code string) error {
	var userAuth mysql.UserAuths
	err := g.GormDB.WithContext(ctx).Model(&userAuth).Where("email = ? AND auth_code = ?", email, code).First(&userAuth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorMsg(ctx, utils.ErrCodeNotFound, utils.Trace(), _errors.ErrCodeNotFound.Error(), utils.ErrFromClient)
		} else {
			return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
		}
	}

	return nil
}

func (g *ValidatePasswordAuthRepository) UpdatePassword(ctx context.Context, user mysql.Users) error {
	err := g.GormDB.WithContext(ctx).Model(&user).Where("email = ?", user.Email).Update("password", &user.Password).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}

	return nil
}

func (g *ValidatePasswordAuthRepository) DeleteAuthCode(ctx context.Context, email string) error {
	userAuth := mysql.UserAuths{}
	err := g.GormDB.WithContext(ctx).Model(&userAuth).Where("email = ?", email).Delete(&userAuth).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}

	return nil
}
