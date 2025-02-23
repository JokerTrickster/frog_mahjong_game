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

type V02GoogleOauthCallbackAuthRepository struct {
	GormDB *gorm.DB
}
type RequestSignupAuthRepository struct {
	GormDB *gorm.DB
}

type FCMTokenAuthRepository struct {
	GormDB *gorm.DB
}

type CheckSigninAuthRepository struct {
	GormDB *gorm.DB
}

type NameCheckAuthRepository struct {
	GormDB *gorm.DB
}

type ValidateSignupAuthRepository struct {
	GormDB *gorm.DB
}