package repository

import (
	_interface "main/features/room/model/interface"

	"gorm.io/gorm"
)

func NewLogoutRoomRepository(gormDB *gorm.DB) _interface.ILogoutRoomRepository {
	return &LogoutRoomRepository{GormDB: gormDB}
}
