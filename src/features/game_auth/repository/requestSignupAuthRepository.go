package repository

import (
	"context"
	_errors "main/features/game_auth/model/errors"
	_interface "main/features/game_auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewRequestSignupAuthRepository(gormDB *gorm.DB) _interface.IRequestSignupAuthRepository {
	return &RequestSignupAuthRepository{GormDB: gormDB}
}

func (g *RequestSignupAuthRepository) FindOneUserByEmail(ctx context.Context, email string) error {
	var user mysql.GameUsers
	result := g.GormDB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.RowsAffected != 0 {
		return utils.ErrorMsg(ctx, utils.ErrUserAlreadyExisted, utils.Trace(), utils.HandleError(_errors.ErrUserAlreadyExisted.Error(), email), utils.ErrFromClient)
	}
	return nil
}

func (d *RequestSignupAuthRepository) InsertAuthCode(ctx context.Context, userAuthDTO mysql.UserAuths) error {

	//인증 코드를 삽입한다.
	err := d.GormDB.Create(&userAuthDTO).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), userAuthDTO), utils.ErrFromMysqlDB)
	}

	return nil
}

func (d *RequestSignupAuthRepository) DeleteAuthCodeByEmail(ctx context.Context, email string) error {
	err := d.GormDB.Where("email = ? and type = ?", email, "signup").Delete(&mysql.UserAuths{}).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), email), utils.ErrFromMysqlDB)
	}
	return nil
}
