package repository

import "gorm.io/gorm"


type ListProfilesRepository struct {
	GormDB *gorm.DB
}

