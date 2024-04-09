package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewWinRequestGameRepository(gormDB *gorm.DB) _interface.IWinRequestGameRepository {
	return &WinRequestGameRepository{GormDB: gormDB}
}

func (d *WinRequestGameRepository) GetRoomUser(ctx context.Context, userID uint, roomID uint) (mysql.RoomUsers, error) {
	var roomUser mysql.RoomUsers
	if err := d.GormDB.Where("user_id = ? AND room_id = ?", userID, roomID).First(&roomUser).Error; err != nil {
		return mysql.RoomUsers{}, utils.ErrorMsg(ctx, utils.ErrBadRequest, utils.Trace(), err.Error(), utils.ErrFromClient)
	}
	return roomUser, nil
}
