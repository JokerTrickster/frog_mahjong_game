package repository

import (
	"context"
	_errors "main/features/game/model/errors"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewDoraGameRepository(gormDB *gorm.DB) _interface.IDoraGameRepository {
	return &DoraGameRepository{GormDB: gormDB}
}

func (d *DoraGameRepository) CheckFirstPlayer(c context.Context, userID int, roomID int) error {
	var roomUsers mysql.RoomUsers
	err := d.GormDB.Model(&roomUsers).Where("user_id = ? AND room_id = ? and turn_number = ?", userID, roomID, 1).First(&roomUsers)
	if err.Error != nil {
		return utils.ErrorMsg(c, utils.ErrBadParameter, utils.Trace(), _errors.ErrNotFirstPlayer.Error(), utils.ErrFromClient)
	}

	return nil
}

func (d *DoraGameRepository) UpdateDoraCard(c context.Context, req *request.ReqDora) error {

	err := d.GormDB.Model(&mysql.Cards{}).Where("id = ? and room_id = ?", req.CardID, req.RoomID).Update("state", "dora")
	if err.Error != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error.Error(), utils.ErrFromMysqlDB)
	}

	return nil
}
