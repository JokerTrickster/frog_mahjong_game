package repository

import "gorm.io/gorm"

type StartGameRepository struct {
	GormDB *gorm.DB
}

type DoraGameRepository struct {
	GormDB *gorm.DB
}

type OwnershipGameRepository struct {
	GormDB *gorm.DB
}
