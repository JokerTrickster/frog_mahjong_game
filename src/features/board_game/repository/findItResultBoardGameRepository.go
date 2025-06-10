package repository

import (
	"context"

	"gorm.io/gorm"

	_interface "main/features/board_game/model/interface"
	"main/utils/db/mysql"
)

func NewFindItResultBoardGameRepository(gormDB *gorm.DB) _interface.IFindItResultBoardGameRepository {
	return &FindItResultBoardGameRepository{GormDB: gormDB}
}

func (d *FindItResultBoardGameRepository) FindFindItResult(ctx context.Context, userID, roomID int) ([]*mysql.FindItUserCorrectPositions, error) {
	var findItResult []*mysql.FindItUserCorrectPositions
	err := d.GormDB.WithContext(ctx).Where("user_id = ? AND room_id = ?", userID, roomID).Find(&findItResult).Error
	return findItResult, err
}

func (d *FindItResultBoardGameRepository) FindGameRoomUser(ctx context.Context, roomID int) ([]*mysql.GameRoomUsers, error) {
	var users []*mysql.GameRoomUsers
	err := d.GormDB.WithContext(ctx).Where("room_id = ?", roomID).Find(&users).Error
	return users, err
}
