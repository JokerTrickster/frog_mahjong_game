package repository

import (
	"context"
	_errors "main/features/game_auth/model/errors"
	_interface "main/features/game_auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewSigninAuthRepository(gormDB *gorm.DB) _interface.ISigninAuthRepository {
	return &SigninAuthRepository{GormDB: gormDB}
}

func (g *SigninAuthRepository) DeleteToken(ctx context.Context, uID uint) error {
	token := mysql.Tokens{
		UserID: uID,
	}
	result := g.GormDB.Model(&token).Where("user_id = ?", uID).Delete(&token)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(result.Error.Error(), uID), utils.ErrFromInternal)
	}
	return nil
}
func (g *SigninAuthRepository) SaveToken(ctx context.Context, uID uint, accessToken, refreshToken string, refreshTknExpiredAt int64) error {
	token := mysql.Tokens{
		UserID:           uID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		RefreshExpiredAt: refreshTknExpiredAt,
	}
	result := g.GormDB.Model(&token).Create(&token)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(result.Error.Error(), uID), utils.ErrFromInternal)
	}
	return nil
}

func (g *SigninAuthRepository) FindOneAndUpdateUser(ctx context.Context, email, password string) (mysql.GameUsers, error) {
	user := mysql.GameUsers{
		Email:  email,
		State:  "wait",
		RoomID: 1,
	}
	// 사용자 정보를 가져온다.
	var findUser mysql.GameUsers
	err := g.GormDB.WithContext(ctx).Model(&findUser).Where("email = ? ", email).First(&findUser).Error
	if err != nil {
		return mysql.GameUsers{}, utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), utils.HandleError(_errors.ErrUserNotFound.Error()+err.Error(), email, password), utils.ErrFromClient)
	}
	if password != findUser.Password {
		return mysql.GameUsers{}, utils.ErrorMsg(ctx, utils.ErrPasswordNotMatch, utils.Trace(), utils.HandleError(_errors.ErrPasswordNotMatch.Error(), password), utils.ErrFromClient)
	}

	result := g.GormDB.WithContext(ctx).Model(&user).Where("email = ?", email).Updates(user)
	if result.Error != nil {
		return mysql.GameUsers{}, utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), utils.HandleError(_errors.ErrUserNotFound.Error()+result.Error.Error(), email, password), utils.ErrFromClient)
	}

	// 변경된 사용자 정보를 가져옵니다.
	err = g.GormDB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return mysql.GameUsers{}, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(err.Error(), email, password), utils.ErrFromInternal)
	}
	return user, nil
}

func (d *SigninAuthRepository) CheckToken(ctx context.Context, uID uint) (*mysql.Tokens, error) {
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

func (g *SigninAuthRepository) DeleteGameRooms(ctx context.Context, uID uint) error {
	result := g.GormDB.WithContext(ctx).Where("owner_id = ?", uID).Delete(&mysql.GameRooms{})
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(result.Error.Error(), uID), utils.ErrFromMysqlDB)
	}
	return nil
}
