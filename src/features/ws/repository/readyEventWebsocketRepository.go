package repository

import (
	"context"
	"errors"
	"log"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"
)

func ReadyFindAllRoomUsers(ctx context.Context, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := mysql.GormMysqlDB.Preload("User").Preload("Room").Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		log.Fatalf("RoomUsers 조회 에러: %s", err)
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
		return errors.New("플레이어 상태 변경 실패")
	}
	return nil
}
