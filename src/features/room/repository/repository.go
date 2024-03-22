package repository

import "gorm.io/gorm"

type CreateRoomRepository struct {
	GormDB *gorm.DB
}

type JoinRoomRepository struct {
	GormDB *gorm.DB
}
