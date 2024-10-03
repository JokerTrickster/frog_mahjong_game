package repository

import (
	"context"
	"main/features/auth/model/entity"
	_errors "main/features/auth/model/errors"
	_interface "main/features/auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewV02GoogleOauthCallbackAuthRepository(gormDB *gorm.DB) _interface.IV02GoogleOauthCallbackAuthRepository {
	return &V02GoogleOauthCallbackAuthRepository{GormDB: gormDB}
}

func (g *V02GoogleOauthCallbackAuthRepository) DeleteToken(ctx context.Context, uID uint) error {
	token := mysql.Tokens{
		UserID: uID,
	}
	result := g.GormDB.Model(&token).Where("user_id = ?", uID).Delete(&token)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(result.Error.Error(),uID), utils.ErrFromInternal)
	}
	return nil
}
func (g *V02GoogleOauthCallbackAuthRepository) SaveToken(ctx context.Context, uID uint, accessToken, refreshToken string, refreshTknExpiredAt int64) error {
	token := mysql.Tokens{
		UserID:           uID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		RefreshExpiredAt: refreshTknExpiredAt,
	}
	result := g.GormDB.Model(&token).Create(&token)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(result.Error.Error(),uID), utils.ErrFromInternal)
	}
	return nil
}

func (g *V02GoogleOauthCallbackAuthRepository) FindOneAndUpdateUser(ctx context.Context, entity *entity.V02GoogleOauthCallbackSQLQuery) (*mysql.Users, error) {
	user := &mysql.Users{
		Email:  entity.Email,
		State:  "wait",
		RoomID: 1,
	}
	//state = "logout"인 유저 wait으로 변경하고 roomID = 1로 변경 user 객체에 반환
	result := g.GormDB.WithContext(ctx).Model(&user).Where("email = ?  ", entity.Email).Updates(&user)
	if result.Error != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), utils.HandleError(_errors.ErrUserNotFound.Error(),entity), utils.ErrFromClient)
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	// 변경된 사용자 정보를 가져옵니다.
	err := g.GormDB.WithContext(ctx).Where("email = ?", entity.Email).First(&user).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(err.Error(),entity), utils.ErrFromInternal)
	}
	return user, nil
}

func (g *V02GoogleOauthCallbackAuthRepository) CreateUser(ctx context.Context, user *mysql.Users) (*mysql.Users, error) {
	result := g.GormDB.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(result.Error.Error(),user), utils.ErrFromInternal)
	}
	return user, nil
}


func (g *V02GoogleOauthCallbackAuthRepository) FindAllBasicProfile(ctx context.Context) ([]*mysql.Profiles, error) {
	profiles := make([]*mysql.Profiles, 0)
	err := g.GormDB.WithContext(ctx).Where("total_count = ?", 0).Find(&profiles).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMysqlDB)
	}
	return profiles, nil
}

func (g *V02GoogleOauthCallbackAuthRepository) InsertOneUserProfile(ctx context.Context, userProfileDTOList []*mysql.UserProfiles) error {
	result := g.GormDB.WithContext(ctx).Create(&userProfileDTOList)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError("failed user profile insert", userProfileDTOList), utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(result.Error.Error(), userProfileDTOList), utils.ErrFromMysqlDB)
	}
	return nil
}
