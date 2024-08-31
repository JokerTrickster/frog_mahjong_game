package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"
)

func ReadyFindAllRoomUsers(ctx context.Context, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := mysql.GormMysqlDB.Preload("User").Preload("Room").Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err)

	}
	return roomUsers, nil
}
func ReadyFindOneAndUpdateRoomUser(ctx context.Context, uID, roomID uint) error {
	// Rooms user에 player state 를 변경한다.
	RoomUser := mysql.RoomUsers{
		PlayerState: "ready",
	}
	err := mysql.GormMysqlDB.Model(&RoomUser).Where("user_id = ? AND room_id = ?", uID, roomID).Updates(RoomUser).Error
	if err != nil {
		return fmt.Errorf("방 유저 정보 업데이트 실패: %v", err)
	}
	return nil
}
