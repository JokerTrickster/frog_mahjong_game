package repository

import (
	"context"
	"fmt"
	"main/features/frog/model/entity"
	"main/features/frog/model/request"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func JoinFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Table("frog_room_users").Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ?", roomID).
		Preload("User").
		Preload("Room").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 실패: %v", err.Error())
	}
	return roomUsers, nil
}

func JoinFindOneRoom(ctx context.Context, tx *gorm.DB, req *request.ReqWSJoin) (mysql.GameRooms, error) {
	// 방 참여 가능한지 체크
	RoomDTO := mysql.GameRooms{}
	result := tx.WithContext(ctx).Where("id = ?", req.RoomID).First(&RoomDTO)
	if result.Error != nil {
		return mysql.GameRooms{}, fmt.Errorf("방 정보를 찾을 수 없습니다. %v", result.Error)
	}
	return RoomDTO, nil
}

func JoinFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) error {
	result := tx.WithContext(ctx).Model(&mysql.GameRooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count + 1"))
	if result.Error != nil {
		return fmt.Errorf("방 인원수 업데이트 실패: %v", result.Error)
	}
	return nil
}

func JoinFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error {
	user := mysql.GameUsers{
		RoomID: int(RoomID),
		State:  "play",
	}
	result := tx.WithContext(ctx).Model(&user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("유저 정보 업데이트 실패: %v", result.Error)
	}

	return nil
}

func JoinInsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomUserDTO mysql.FrogRoomUsers) error {
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
	result := tx.WithContext(ctx).Where("user_id = ? and room_id = ?", uID, roomID).Delete(&mysql.FrogRoomUsers{})
	// 방 유저 정보가 없는 경우
	if result.RowsAffected == 0 {
		return nil
	} else {
		result2 := tx.WithContext(ctx).Model(&mysql.GameRooms{}).Where("id = ?", roomID).Update("current_count", gorm.Expr("current_count - 1"))
		if result2.Error != nil {
			return fmt.Errorf("방 인원수 업데이트 실패: %v", result.Error)
		}
	}
	if result.Error != nil {
		return fmt.Errorf("failed to delete room user: %v", result.Error)
	}
	return nil
}
