package repository

import (
	_interface "main/features/room/model/interface"

	"gorm.io/gorm"
)

func NewListRoomRepository(gormDB *gorm.DB) _interface.IListRoomRepository {
	return &ListRoomRepository{GormDB: gormDB}
}
