package repository

import (
	"context"

	"gorm.io/gorm"

	_interface "main/features/board_game/model/interface"
	"main/utils/db/mysql"
)

func NewSequenceResultBoardGameRepository(gormDB *gorm.DB) _interface.ISequenceResultBoardGameRepository {
	return &SequenceResultBoardGameRepository{GormDB: gormDB}
}

func (d *SequenceResultBoardGameRepository) FindGameRoomUser(ctx context.Context, roomID int) ([]*mysql.SequenceUsers, error) {
	var users []*mysql.SequenceUsers
	err := d.GormDB.WithContext(ctx).Where("room_id = ?", roomID).Find(&users).Error
	return users, err
}

func (d *SequenceResultBoardGameRepository) FindRoomMaps(ctx context.Context, roomID int) ([]*mysql.SequenceRoomMaps, error) {
	var maps []*mysql.SequenceRoomMaps
	err := d.GormDB.WithContext(ctx).Where("room_id = ?", roomID).Find(&maps).Error
	return maps, err
}

func (d *SequenceResultBoardGameRepository) FindGameResult(ctx context.Context, roomID int) ([]*mysql.GameResults, error) {
	var results []*mysql.GameResults
	err := d.GormDB.WithContext(ctx).Where("room_id = ?", roomID).Find(&results).Error
	return results, err
}
