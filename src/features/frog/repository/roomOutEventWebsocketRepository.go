package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func RoomOutFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}

func RoomOutCheckOwner(ctx context.Context, tx *gorm.DB, uID uint, roomID uint) error {
	// 방장인지 체크
	room := mysql.Rooms{}
	err := tx.WithContext(ctx).Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return fmt.Errorf("방 정보를 찾을 수 없습니다. %v", err)
	}
	if room.OwnerID != int(uID) {
		return fmt.Errorf("방장만 강제퇴장할 수 있습니다.")
	}
	return nil
}

func RoomOutUpdateUser(ctx context.Context, tx *gorm.DB, targetUserID uint, roomID uint) error {
	// 타겟 유저 데이터 변경 (플레이 상태, 룸ID)
	err := tx.WithContext(ctx).Model(&mysql.Users{}).Where("id = ?", targetUserID).Updates(mysql.Users{RoomID: int(1), State: "wait"})
	if err.Error != nil {
		return fmt.Errorf("유저 정보 업데이트 실패: %v", err.Error)
	}
	return nil
}

func RoomOutDeleteRoomUser(ctx context.Context, tx *gorm.DB, targetUserID uint, roomID uint) error {
	// 룸 유저 정보 삭제
	err := tx.WithContext(ctx).Model(&mysql.FrogRoomUsers{}).Where("user_id = ? and room_id = ?", targetUserID, roomID).Delete(&mysql.FrogRoomUsers{})

	if err.Error != nil {
		return fmt.Errorf("룸 유저 정보 삭제 실패 %v", err.Error)
	}
	return nil
}

func RoomOutUpdateRoom(ctx context.Context, tx *gorm.DB, roomID uint) error {
	// 방 현재 인원을 감소시킨다.
	err := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", roomID).Update("current_count", gorm.Expr("current_count - 1"))
	if err.Error != nil {
		return fmt.Errorf("방 인원수 업데이트 실패: %v", err.Error)
	}
	return nil
}
