package repository

import (
	"gorm.io/gorm"
)

type SignupAuthRepository struct {
	GormDB *gorm.DB
}
type SigninAuthRepository struct {
	GormDB *gorm.DB
}

type LogoutAuthRepository struct {
	GormDB *gorm.DB
}

type ReissueAuthRepository struct {
	GormDB *gorm.DB
}

type GoogleOauthAuthRepository struct {
	GormDB *gorm.DB
}
type GoogleOauthCallbackAuthRepository struct {
	GormDB *gorm.DB
}

type RequestPasswordAuthRepository struct {
	GormDB *gorm.DB
}

type ValidatePasswordAuthRepository struct {
	GormDB *gorm.DB
}