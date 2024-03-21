package repository

import (
	"context"
	"fmt"
	_interface "main/features/room/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewCreateRoomRepository(gormDB *gorm.DB) _interface.ICreateRoomRepository {
	return &CreateRoomRepository{GormDB: gormDB}
}
func (g *CreateRoomRepository) InsertOneRoom(ctx context.Context, roomDTO mysql.Rooms) (int, error) {
	result := g.GormDB.WithContext(ctx).Create(&roomDTO)
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("failed room insert one")
	}
	if result.Error != nil {
		return 0, result.Error
	}
	return int(roomDTO.ID), nil
}
func (g *CreateRoomRepository) InsertOneRoomUser(ctx context.Context, roomUserDTO mysql.RoomUsers) error {
	result := g.GormDB.WithContext(ctx).Create(&roomUserDTO)
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed room user insert one")
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
