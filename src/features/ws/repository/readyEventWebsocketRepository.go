package repository

import (
	"context"
	"errors"
	"main/utils/db/mysql"
)

func ReadyFindOneAndUpdateRoomUser(ctx context.Context, uID, roomID uint) error {
	// Rooms user에 player state 를 변경한다.
	RoomUser := mysql.RoomUsers{
		PlayerState: "ready",
	}
	err := mysql.GormMysqlDB.Model(&RoomUser).Where("user_id = ? AND room_id = ?", uID, roomID).Updates(RoomUser).Error
	if err != nil {
		return errors.New("플레이어 상태 변경 실패")
	}
	return nil
}
