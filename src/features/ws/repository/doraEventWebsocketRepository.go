package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func DoraFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
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
func DoraCheckFirstPlayer(ctx context.Context, tx *gorm.DB, userID, roomID uint) error {
	var roomUser mysql.FrogRoomUsers
	err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("user_id = ?", userID).
		Where("room_id = ?", roomID).
		Where("turn_number = ?", 1).
		First(&roomUser).Error
	if err != nil {
		return fmt.Errorf("첫 번째 플레이어가 아닙니다: %v", err)
	}
	return nil
}

func DoraUpdateDoraCard(ctx context.Context, tx *gorm.DB, entity *entity.WSDoraEntity) error {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", entity.RoomID).
		Where("card_id = ?", entity.CardID).
		Update("state", "dora").Error
	if err != nil {
		return fmt.Errorf("도라 카드 업데이트 실패: %v", err)
	}
	return nil
}
