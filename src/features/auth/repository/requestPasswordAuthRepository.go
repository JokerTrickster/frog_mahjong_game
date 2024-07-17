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

func NewRequestPasswordAuthRepository(gormDB *gorm.DB) _interface.IRequestPasswordAuthRepository {
	return &RequestPasswordAuthRepository{GormDB: gormDB}
}

func (g *RequestPasswordAuthRepository) FindOneUserByEmail(ctx context.Context, email string) error {
	var user mysql.Users
	result := g.GormDB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), _errors.ErrUserNotFound.Error(), utils.ErrFromClient)
	} 
	return nil
}

// 이메일로 찾아서 있으면 업데이트하고 없으면 삽입한다.
func (d *RequestPasswordAuthRepository) InsertAuthCode(ctx context.Context, userAuthDTO mysql.UserAuths) error {

	var existingUserAuth mysql.UserAuths

	// 이메일로 기존 레코드가 있는지 확인
	err := d.GormDB.Where("email = ?", userAuthDTO.Email).First(&existingUserAuth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 레코드가 없으면 삽입
			err = d.GormDB.Create(&userAuthDTO).Error
			if err != nil {
				return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
			}
		} else {
			// 다른 에러가 발생하면 에러 반환
			return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
		}
	} else {
		// 레코드가 있으면 업데이트
		userAuthDTO.ID = existingUserAuth.ID
		err = d.GormDB.WithContext(ctx).Model(&userAuthDTO).Where("email = ?", userAuthDTO.Email).Update("auth_code", &userAuthDTO.AuthCode).Error
		if err != nil {
			return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
		}
	}

	return nil
}
