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

func NewJoinPlayRoomsRepository(gormDB *gorm.DB) _interface.IJoinPlayRoomsRepository {
	return &JoinPlayRoomsRepository{GormDB: gormDB}
}

func (g *JoinPlayRoomsRepository) FindOneRoom(ctx context.Context, req *request.ReqJoinPlay) error {
	// 방 참여 가능한지 체크
	RoomDTO := mysql.Rooms{}
	err := mysql.GormMysqlDB.WithContext(ctx).Where("deleted_at is null and password = ? and state = ?", req.Password, "wait").First(&RoomDTO).Error
	if err != nil {
		if err.Error() == "record not found" {
			return utils.ErrorMsg(ctx, utils.ErrWrongPassword, utils.Trace(), _errors.ErrWrongPassword.Error(), utils.ErrFromClient)
		}
		return utils.ErrorMsg(ctx, utils.ErrRoomNotFound, utils.Trace(), _errors.ErrRoomNotFound.Error(), utils.ErrFromClient)
	}
	if RoomDTO.CurrentCount == RoomDTO.MaxCount {
		return utils.ErrorMsg(ctx, utils.ErrRoomFull, utils.Trace(), _errors.ErrRoomFull.Error(), utils.ErrFromClient)
	}
	return nil
}
