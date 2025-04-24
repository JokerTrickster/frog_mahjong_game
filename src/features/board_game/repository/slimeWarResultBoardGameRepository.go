package repository

import (
	"context"

	"gorm.io/gorm"

	_interface "main/features/board_game/model/interface"
	"main/utils/db/mysql"
)

func NewSlimeWarResultBoardGameRepository(gormDB *gorm.DB) _interface.ISlimeWarResultBoardGameRepository {
	return &SlimeWarResultBoardGameRepository{GormDB: gormDB}
}

func (d *SlimeWarResultBoardGameRepository) FindGameRoomUser(ctx context.Context, roomID int) ([]*mysql.SlimeWarUsers, error) {
	var users []*mysql.SlimeWarUsers
	err := d.GormDB.WithContext(ctx).Where("room_id = ?", roomID).Find(&users).Error
	return users, err
}

func (d *SlimeWarResultBoardGameRepository) FindRoomMaps(ctx context.Context, roomID int) ([]*mysql.SlimeWarRoomMaps, error) {
	var maps []*mysql.SlimeWarRoomMaps
	err := d.GormDB.WithContext(ctx).Where("room_id = ?", roomID).Find(&maps).Error
	return maps, err
}
