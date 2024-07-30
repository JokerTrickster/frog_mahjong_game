package repository

import (
	_interface "main/features/rooms/model/interface"

	"gorm.io/gorm"
)

func NewV02JoinRoomsRepository(gormDB *gorm.DB) _interface.IV02JoinRoomsRepository {
	return &V02JoinRoomsRepository{GormDB: gormDB}
}
