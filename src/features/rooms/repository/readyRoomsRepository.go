package repository

import (
	"context"
	_errors "main/features/rooms/model/errors"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewReadyRoomsRepository(gormDB *gorm.DB) _interface.IReadyRoomsRepository {
	return &ReadyRoomsRepository{GormDB: gormDB}
}

func (g *ReadyRoomsRepository) FindOneAndUpdateRoomUser(ctx context.Context, uID uint, req *request.ReqReady) error {
	// Rooms user에 player state 를 변경한다.
	RoomUser := mysql.RoomUsers{
		PlayerState: req.PlayerState,
	}
	err := g.GormDB.Model(&RoomUser).Where("user_id = ? AND Room_id = ?", uID, req.RoomID).Updates(RoomUser).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), _errors.ErrPlayerStateFailed.Error(), utils.ErrFromClient)
	}
	return nil
}
