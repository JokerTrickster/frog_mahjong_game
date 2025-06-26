package repository

import "gorm.io/gorm"

type ListProfilesRepository struct {
	GormDB *gorm.DB
}

type UploadProfilesRepository struct {
	GormDB *gorm.DB
}

type UpdateProfilesRepository struct {
	GormDB *gorm.DB
}
