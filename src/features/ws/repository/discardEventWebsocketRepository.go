package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func DiscardCardsFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
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
func DiscardCardsFindOneDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, error) {
	var dora mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("state = ?", "dora").
		First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다: %v", err)
	}
	return &dora, nil
}

func DiscardCardsUpdateCardState(ctx context.Context, tx *gorm.DB, entity *entity.WSDiscardCardsEntity) error {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", entity.RoomID).
		Where("card_id = ?", entity.CardID).
		Where("state = ?", "owned").
		Updates(map[string]interface{}{
			"state":   "discard",
			"user_id": int(entity.UserID),
		}).Error
	if err != nil {
		return fmt.Errorf("카드 상태 업데이트 실패: %v", err)
	}
	return nil
}

func DiscardCardsUpdateRoomUserCardCount(ctx context.Context, tx *gorm.DB, entity *entity.WSDiscardCardsEntity) error {
	err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("room_id = ?", entity.RoomID).
		Where("user_id = ?", entity.UserID).
		Update("owned_card_count", gorm.Expr("owned_card_count - 1")).Error
	if err != nil {
		return fmt.Errorf("방 유저 카드 카운트 업데이트 실패: %v", err)
	}
	return nil
}

func DiscardCardsFindAllCard(ctx context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.FrogUserCards, error) {
	var cards []*mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("user_id = ?", userID).
		Order("updated_at ASC").
		Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다: %v", err)
	}
	return cards, nil
}
