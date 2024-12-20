package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func AbnormalFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("RoomUsers").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}

// 카드 정보 모두 삭제
func AbnormalDeleteAllCards(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) error {
	err := tx.Model(&mysql.Cards{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.Cards{}).Error
	if err != nil {
		return fmt.Errorf("카드 삭제 실패 %v", err.Error())
	}
	return nil
}

// 방 삭제 처리
func AbnormalDeleteRoom(c context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) error {
	err := tx.Model(&mysql.Rooms{}).Where("id = ?", AbnormalEntity.RoomID).Delete(&mysql.Rooms{}).Error
	if err != nil {
		return fmt.Errorf("방 삭제 실패 %v", err.Error())
	}
	return nil
}

// 룸 유저 상태 변경 (play -> wait)
func AbnormalUpdateRoomUsers(c context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) error {
	err := tx.Model(&mysql.RoomUsers{}).Where("room_id = ? and user_id = ?", AbnormalEntity.RoomID, AbnormalEntity.UserID).Update("player_state", "disconnected").Error
	if err != nil {
		return fmt.Errorf("유저 상태 변경 실패 %v", err.Error())
	}
	return nil
}
