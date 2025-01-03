package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TimeOutDiscardCardsFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
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
func TimeOutDiscardFindDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, error) {
	var dora mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ? AND state = ?", roomID, "dora").
		First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다: %v", err.Error())
	}
	return &dora, nil
}

func TimeOutDiscardUpdateCardState(ctx context.Context, tx *gorm.DB, entity *entity.WSTimeOutDiscardCardsEntity) error {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ? AND card_id = ? AND state = ?", entity.RoomID, entity.CardID, "none").
		Updates(&mysql.FrogUserCards{
			State:  "discard",
			UserID: int(entity.UserID),
		}).Error
	if err != nil {
		return fmt.Errorf("카드 상태 업데이트 실패: %v", err.Error())
	}
	return nil
}

func TimeOutDiscardFindAllCards(ctx context.Context, tx *gorm.DB, roomID, userID uint) ([]*mysql.FrogUserCards, error) {
	var cards []*mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Order("updated_at ASC").
		Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다: %v", err.Error())
	}
	return cards, nil
}
