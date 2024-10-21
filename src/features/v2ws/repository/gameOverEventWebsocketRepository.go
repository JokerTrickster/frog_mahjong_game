package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func GameOverFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}

// 카드 정보 모두 삭제
func GameOverDeleteAllCards(ctx context.Context, tx *gorm.DB, GameOverEntity *entity.WSGameOverEntity) error {
	err := tx.Model(&mysql.Cards{}).Where("room_id = ?", GameOverEntity.RoomID).Delete(&mysql.Cards{}).Error
	if err != nil {
		return fmt.Errorf("카드 삭제 실패 %v", err.Error())
	}
	return nil
}

// 유저 상태 변경 (play -> wait)
func GameOverUpdateRoomUsers(c context.Context, tx *gorm.DB, GameOverEntity *entity.WSGameOverEntity) error {
	err := tx.Model(&mysql.RoomUsers{}).Where("room_id = ?", GameOverEntity.RoomID).Update("player_state", "wait").Error
	if err != nil {
		return fmt.Errorf("방 유저 상태 변경 실패 %v", err.Error())
	}
	return nil
}
