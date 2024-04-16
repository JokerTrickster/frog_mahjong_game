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
