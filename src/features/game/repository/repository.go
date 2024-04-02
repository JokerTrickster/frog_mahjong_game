package repository

import "gorm.io/gorm"

type StartGameRepository struct {
	GormDB *gorm.DB
}
