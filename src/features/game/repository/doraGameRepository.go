package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewDoraGameRepository(gormDB *gorm.DB) _interface.IDoraGameRepository {
	return &DoraGameRepository{GormDB: gormDB}
}

func (d *DoraGameRepository) CheckOwner(c context.Context, userID int, roomID int) error {
	var roomUsers mysql.RoomUsers
	err := d.GormDB.Model(&roomUsers).Where("user_id = ? AND room_id = ? and player_state = ?", userID, roomID).First(&roomUsers)

	if err != nil {
		return err.Error
	}

	return nil
}

func (d *DoraGameRepository) UpdateDoraCard(c context.Context, req *request.ReqDora) error {
	return nil
}
