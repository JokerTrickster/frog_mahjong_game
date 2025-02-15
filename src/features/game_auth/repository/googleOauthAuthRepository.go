package repository

import (
	_interface "main/features/game_auth/model/interface"

	"gorm.io/gorm"
)

func NewGoogleOauthAuthRepository(gormDB *gorm.DB) _interface.IGoogleOauthAuthRepository {
	return &GoogleOauthAuthRepository{GormDB: gormDB}
}
