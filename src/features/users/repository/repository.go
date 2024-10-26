package repository

import "gorm.io/gorm"

type GetUsersRepository struct {
	GormDB *gorm.DB
}

type ListUsersRepository struct {
	GormDB *gorm.DB
}

type UpdateUsersRepository struct {
	GormDB *gorm.DB
}
type DeleteUsersRepository struct {
	GormDB *gorm.DB
}

type ListProfilesUsersRepository struct {
	GormDB *gorm.DB
}

type FullCoinUsersRepository struct {
	GormDB *gorm.DB
}
