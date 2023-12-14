package repository

import (
	"gorm.io/gorm"
)

type SignupAuthRepository struct {
	GormDB *gorm.DB
}
