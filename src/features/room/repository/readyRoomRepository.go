package repository

import (
	_interface "main/features/room/model/interface"

	"gorm.io/gorm"
)

func NewReadyRoomRepository(gormDB *gorm.DB) _interface.IReadyRoomRepository {
	return &ReadyRoomRepository{GormDB: gormDB}
}
