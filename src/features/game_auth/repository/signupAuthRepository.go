package repository

import (
	"context"
	_errors "main/features/game_auth/model/errors"
	_interface "main/features/game_auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewSignupAuthRepository(gormDB *gorm.DB) _interface.ISignupAuthRepository {
	return &SignupAuthRepository{GormDB: gormDB}
}
func (g *SignupAuthRepository) UserCheckByEmail(ctx context.Context, email string) error {
	var user mysql.GameUsers
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
func (g *SignupAuthRepository) InsertOneUser(ctx context.Context, user *mysql.GameUsers) error {
	result := g.GormDB.WithContext(ctx).Create(&user)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError("failed user insert", user), utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(result.Error.Error(), user), utils.ErrFromMysqlDB)
	}
	return nil
}

func (g *SignupAuthRepository) VerifyAuthCode(ctx context.Context, email, code string) error {
	var userAuth mysql.UserAuths

	result := g.GormDB.WithContext(ctx).Where("email = ? AND auth_code = ? and type = ? and project = ? and is_active = ?", email, code, "signup", "board_game", true).First(&userAuth)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInvalidAuthCode, utils.Trace(), utils.HandleError(_errors.ErrInvalidAuthCode.Error(), email, code), utils.ErrFromClient)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(result.Error.Error(), email, code), utils.ErrFromMysqlDB)
	}
	return nil
}

func (g *SignupAuthRepository) FindAllBasicProfile(ctx context.Context) ([]*mysql.GameProfiles, error) {
	profiles := make([]*mysql.GameProfiles, 0)
	err := g.GormDB.WithContext(ctx).Find(&profiles).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMysqlDB)
	}
	return profiles, nil
}

func (g *SignupAuthRepository) InsertOneUserProfile(ctx context.Context, userProfileDTOList []*mysql.GameUserProfiles) error {

	result := g.GormDB.WithContext(ctx).Create(&userProfileDTOList)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError("failed user profile insert", userProfileDTOList), utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(result.Error.Error(), userProfileDTOList), utils.ErrFromMysqlDB)
	}
	return nil
}
