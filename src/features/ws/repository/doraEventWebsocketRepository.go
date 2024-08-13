package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func DoraFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}
func DoraCheckFirstPlayer(c context.Context, tx *gorm.DB, userID uint, roomID uint) error {
	var roomUsers mysql.RoomUsers
	err := tx.Model(&roomUsers).Where("user_id = ? AND room_id = ? and turn_number = ?", userID, roomID, 1).First(&roomUsers)
	if err.Error != nil {
		return fmt.Errorf("첫번째 플레이어가 아닙니다. %v", err.Error)
	}

	return nil
}

func DoraUpdateDoraCard(c context.Context, tx *gorm.DB, entity *entity.WSDoraEntity) error {

	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and card_id = ?", entity.RoomID, entity.CardID).Update("state", "dora")
	if err.Error != nil {
		return fmt.Errorf("도라 카드 업데이트 실패 %v", err.Error)
	}

	return nil
}
