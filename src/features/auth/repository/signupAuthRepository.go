package repository

import (
	"context"
	_errors "main/features/auth/model/errors"
	_interface "main/features/auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewSignupAuthRepository(gormDB *gorm.DB) _interface.ISignupAuthRepository {
	return &SignupAuthRepository{GormDB: gormDB}
}
func (g *SignupAuthRepository) UserCheckByEmail(ctx context.Context, email string) error {
	var user mysql.Users
	result := g.GormDB.WithContext(ctx).Where("email = ? ", email).First(&user)
	if result.RowsAffected == 0 {
		return nil
	} else {
		if user.Provider == "email" {
			return utils.ErrorMsg(ctx, utils.ErrUserAlreadyExisted, utils.Trace(), utils.HandleError(_errors.ErrUserAlreadyExisted.Error(), email), utils.ErrFromClient)

		} else {
			//google 이미 가입 된 상태
			return utils.ErrorMsg(ctx, utils.ErrUserGoogleAlreadyExisted, utils.Trace(), utils.HandleError(_errors.ErrUserGoogleExisted.Error(), email), utils.ErrFromClient)
		}
	}
}
func (g *SignupAuthRepository) InsertOneUser(ctx context.Context, user mysql.Users) error {
	result := g.GormDB.WithContext(ctx).Create(&user)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError("failed user insert", user), utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(result.Error.Error(), user), utils.ErrFromMysqlDB)
	}
	return nil
}
