package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ImportCardsFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
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
func ImportCardsFindOneDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, error) {
	var dora mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("state = ?", "dora").
		First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다: %v", err.Error())
	}
	return &dora, nil
}
func ImportCardsUpdateCardState(c context.Context, tx *gorm.DB, entity *entity.WSImportCardsEntity) error {
	for _, card := range entity.Cards {
		err := tx.Model(&mysql.FrogUserCards{}).
			Where("room_id = ?", card.RoomID).
			Where("card_id = ?", card.CardID).
			Where("state = ?", "none").
			Updates(&mysql.FrogUserCards{
				State:  "owned",
				UserID: card.UserID,
			}).Error
		if err != nil {
			return fmt.Errorf("카드 상태 업데이트 실패: %v", err.Error())
		}
	}
	return nil
}

func ImportCardsUpdateRoomUserCardCount(ctx context.Context, tx *gorm.DB, entity *entity.WSImportCardsEntity) error {
	// 유저 ID와 Room ID를 기준으로 owned_card_count를 증가시킵니다.
	for _, card := range entity.Cards {
		err := tx.Model(&mysql.FrogRoomUsers{}).
			Where("room_id = ?", card.RoomID).
			Where("user_id = ?", card.UserID).
			Update("owned_card_count", gorm.Expr("owned_card_count + 1")).Error
		if err != nil {
			return fmt.Errorf("방 유저 카드 카운트 업데이트 실패: %v", err.Error())
		}
	}
	return nil
}
func ImportCardsFindAllCard(ctx context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.FrogUserCards, error) {
	var cards []*mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("user_id = ?", userID).
		Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다: %v", err.Error())
	}
	return cards, nil
}
