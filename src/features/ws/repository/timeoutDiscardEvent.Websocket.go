package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func TimeOutDiscardCardsFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}
func TimeOutDiscardCardsFindOneDora(c context.Context, tx *gorm.DB, roomID uint) (*mysql.Cards, error) {
	dora := mysql.Cards{}
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and state = ?", roomID, "dora").First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다. %v", err.Error())
	}
	return &dora, nil
}
func TimeOutDiscardCardsUpdateCardState(c context.Context, tx *gorm.DB, entity *entity.WSTimeOutDiscardCardsEntity) error {
	// 카드 상태 업데이트
	// room_id, card_id, state로 찾고 카드 업데이트할 때 트랜잭션 처리해줘
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and card_id = ? and state = ?", entity.RoomID, entity.CardID, "none").Updates(&mysql.Cards{State: "discard", UserID: int(entity.UserID)}).Error
	if err != nil {
		return fmt.Errorf("카드 버리기 상태 업데이트 실패 %v", err.Error())
	}
	return nil
}

func TimeOutDiscardCardsFindAllCard(c context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.Cards, error) {
	cards := make([]*mysql.Cards, 0)
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and user_id = ?", roomID, userID).Order("updated_at ASC").Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다. %v", err.Error())
	}
	return cards, nil
}
