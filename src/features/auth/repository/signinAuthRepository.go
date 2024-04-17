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
func (g *SigninAuthRepository) DeleteToken(ctx context.Context, uID uint) error {
	token := mysql.Tokens{
		UserID: uID,
	}
	result := g.GormDB.Model(&token).Where("user_id = ?", uID).Delete(&token)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), result.Error.Error(), utils.ErrFromInternal)
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
		return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), result.Error.Error(), utils.ErrFromInternal)
	}
	return nil
}

func (g *SigninAuthRepository) FindOneAndUpdateUser(ctx context.Context, email, password string) (mysql.Users, error) {
	user := mysql.Users{
		Email:  email,
		State:  "wait",
		RoomID: 1,
	}
	//state = "logout"인 유저 wait으로 변경하고 roomID = 1로 변경 user 객체에 반환
	result := g.GormDB.WithContext(ctx).Model(&user).Where("email = ? and password = ? and state = ?", email, password, "logout").Updates(user)
	if result.Error != nil {
		return mysql.Users{}, utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), _errors.ErrUserNotFound.Error(), utils.ErrFromClient)
	}
	if result.RowsAffected == 0 {
		return mysql.Users{}, utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), _errors.ErrUserNotFound.Error(), utils.ErrFromClient)
	}
	// 변경된 사용자 정보를 가져옵니다.
	err := g.GormDB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return mysql.Users{}, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), err.Error(), utils.ErrFromInternal)
	}
	return user, nil
}
