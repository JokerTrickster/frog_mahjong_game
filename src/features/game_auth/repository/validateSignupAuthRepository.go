package repository

import (
	"context"
	"errors"
	_interface "main/features/game_auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"
	"time"

	"gorm.io/gorm"
)

func NewValidateSignupAuthRepository(gormDB *gorm.DB) _interface.IValidateSignupAuthRepository {
	return &ValidateSignupAuthRepository{GormDB: gormDB}
}

// CheckAuthCode
func (g *ValidateSignupAuthRepository) CheckAuthCode(ctx context.Context, userAuthDTO *mysql.UserAuths) error {
	var userAuth mysql.UserAuths

	timeLimit := time.Now().Add(-10 * time.Minute)
	err := g.GormDB.WithContext(ctx).Where("email = ? AND auth_code = ? AND type = ? AND project = ? AND created_at >= ?",
		userAuthDTO.Email, userAuthDTO.AuthCode, userAuthDTO.Type, userAuthDTO.Project, timeLimit).First(&userAuth).Error

	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.ErrorMsg(ctx, utils.ErrInvalidAuthCode, utils.Trace(), "Invalid or expired auth code", utils.ErrFromClient)
	}
	return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), err.Error(), utils.ErrFromInternal)
}
func (g *ValidateSignupAuthRepository) UpdateAuthCode(ctx context.Context, userAuthDTO *mysql.UserAuths) error {
	// is_active를 true로 업데이트
	err := g.GormDB.WithContext(ctx).Model(&mysql.UserAuths{}).
		Where("email = ? AND auth_code = ? AND type = ? AND project = ?", userAuthDTO.Email, userAuthDTO.AuthCode, userAuthDTO.Type, userAuthDTO.Project).
		Update("is_active", true).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), err.Error(), utils.ErrFromInternal)
	}
	return nil
}
