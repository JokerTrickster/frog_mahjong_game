package repository

import "gorm.io/gorm"

type GetUsersRepository struct {
	GormDB *gorm.DB
}