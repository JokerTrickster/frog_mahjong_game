package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/features/ws/model/request"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func JoinFindAllRoomUsers(ctx context.Context, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := mysql.GormMysqlDB.Preload("User").Preload("Room").Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err)
	}
	return roomUsers, nil
}

func JoinFindOneRoom(ctx context.Context, tx *gorm.DB, req *request.ReqWSJoin) (mysql.Rooms, error) {
	// 방 참여 가능한지 체크
	RoomDTO := mysql.Rooms{}
	result := tx.WithContext(ctx).Where("id = ?", req.RoomID).First(&RoomDTO)
	if result.Error != nil {
		return mysql.Rooms{}, fmt.Errorf("방 정보를 찾을 수 없습니다. %v", result.Error)
	}
	return RoomDTO, nil
}

func JoinFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) error {
	result := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count + 1"))
	if result.Error != nil {
		return fmt.Errorf("방 인원수 업데이트 실패: %v", result.Error)
	}
	return nil
}

func JoinFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "play",
	}
	result := tx.WithContext(ctx).Model(&user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("유저 정보 업데이트 실패: %v", result.Error)
	}

	return nil
}

func JoinInsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomUserDTO mysql.RoomUsers) error {
	result := tx.WithContext(ctx).Create(&RoomUserDTO)
	if result.RowsAffected == 0 {
		return fmt.Errorf("방 유저 정보 생성 실패")
	}
	if result.Error != nil {
		return fmt.Errorf("방 유저 정보 생성 실패: %v", result.Error)
	}
	return nil
}

func JoinFindOneAndDeleteRoomUser(ctx context.Context, tx *gorm.DB, uID, roomID uint) error {
	result := tx.WithContext(ctx).Where("user_id = ? and room_id = ?", uID, roomID).Delete(&mysql.RoomUsers{})
	// 방 유저 정보가 없는 경우
	if result.RowsAffected == 0 {
		return nil
	} else {
		result2 := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", roomID).Update("current_count", gorm.Expr("current_count - 1"))
		if result2.Error != nil {
			return fmt.Errorf("방 인원수 업데이트 실패: %v", result.Error)
		}
	}
	if result.Error != nil {
		return fmt.Errorf("failed to delete room user: %v", result.Error)
	}
	return nil
}
