package repository

import (
	"context"
	"errors"
	"main/features/game_auth/model/entity"
	_interface "main/features/game_auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewGoogleOauthCallbackAuthRepository(gormDB *gorm.DB) _interface.IGoogleOauthCallbackAuthRepository {
	return &GoogleOauthCallbackAuthRepository{GormDB: gormDB}
}

func (g *GoogleOauthCallbackAuthRepository) DeleteToken(ctx context.Context, uID uint) error {
	token := mysql.Tokens{
		UserID: uID,
	}
	result := g.GormDB.Model(&token).Where("user_id = ?", uID).Delete(&token)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(result.Error.Error(), uID), utils.ErrFromInternal)
	}
	return nil
}
func (g *GoogleOauthCallbackAuthRepository) SaveToken(ctx context.Context, uID uint, accessToken, refreshToken string, refreshTknExpiredAt int64) error {
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

func (g *GoogleOauthCallbackAuthRepository) FindOneAndUpdateUser(ctx context.Context, entity *entity.GoogleOauthCallbackSQLQuery) (*mysql.GameUsers, error) {
	var gameUser mysql.GameUsers

	// 이메일로 사용자를 조회한다.
	err := g.GormDB.WithContext(ctx).Where("email = ?", entity.Email).First(&gameUser).Error
	if err != nil {
		//레코드를 찾지 못한 경우 (nil, nil) 반환
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), entity), utils.ErrFromMysqlDB)
	}

	//업데이트된 사용자 정보를 저장한다.
	err = g.GormDB.WithContext(ctx).Model(&gameUser).Updates(map[string]interface{}{
		"state":   "wait",
		"room_id": 1,
	}).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), entity), utils.ErrFromMysqlDB)
	}

	return &gameUser, nil
}

func (g *GoogleOauthCallbackAuthRepository) CreateUser(ctx context.Context, user *mysql.GameUsers) (*mysql.GameUsers, error) {
	result := g.GormDB.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), utils.HandleError(result.Error.Error(), user), utils.ErrFromInternal)
	}
	return user, nil
}

func (g *GoogleOauthCallbackAuthRepository) FindAllBasicProfile(ctx context.Context) ([]*mysql.GameProfiles, error) {
	profiles := make([]*mysql.GameProfiles, 0)
	err := g.GormDB.WithContext(ctx).Find(&profiles).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMysqlDB)
	}
	return profiles, nil
}

func (g *GoogleOauthCallbackAuthRepository) InsertOneUserProfile(ctx context.Context, userProfileDTOList []*mysql.GameUserProfiles) error {
	result := g.GormDB.WithContext(ctx).Create(&userProfileDTOList)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError("failed user profile insert", userProfileDTOList), utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(result.Error.Error(), userProfileDTOList), utils.ErrFromMysqlDB)
	}
	return nil
}

func (d *GoogleOauthCallbackAuthRepository) CheckToken(ctx context.Context, uID uint) (*mysql.Tokens, error) {
	token := &mysql.Tokens{}
	err := d.GormDB.WithContext(ctx).Where("user_id = ?", uID).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), uID), utils.ErrFromMysqlDB)
	}
	return token, nil
}
