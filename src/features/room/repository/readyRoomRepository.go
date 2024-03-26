package repository

import (
	"context"
	_errors "main/features/room/model/errors"
	_interface "main/features/room/model/interface"
	"main/features/room/model/request"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewReadyRoomRepository(gormDB *gorm.DB) _interface.IReadyRoomRepository {
	return &ReadyRoomRepository{GormDB: gormDB}
}

func (g *ReadyRoomRepository) FindOneAndUpdateRoomUser(ctx context.Context, uID uint, req *request.ReqReady) error {
	// room user에 player state 를 변경한다.
	roomUser := mysql.RoomUsers{
		PlayerState: req.PlayerState,
	}
	err := g.GormDB.Model(&roomUser).Where("user_id = ? AND room_id = ?", uID, req.RoomID).Updates(roomUser).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), _errors.ErrPlayerStateFailed.Error(), utils.ErrFromClient)
	}
	return nil
}
