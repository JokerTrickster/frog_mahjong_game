package repository

import (
	_interface "main/features/room/model/interface"

	"gorm.io/gorm"
)

func NewUserListRoomRepository(gormDB *gorm.DB) _interface.IUserListRoomRepository {
	return &UserListRoomRepository{GormDB: gormDB}
}
