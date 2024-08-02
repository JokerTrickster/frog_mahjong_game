package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func ImportCardsFindAllRoomUsers(ctx context.Context, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := mysql.GormMysqlDB.Preload("User").Preload("Room").Preload("Cards", "room_id = ?", roomID).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}
func ImportCardsFindOneDora(c context.Context, tx *gorm.DB, roomID uint) (*mysql.Cards, error) {
	dora := mysql.Cards{}
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and state = ?", roomID, "dora").First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다. %v", err.Error())
	}
	return &dora, nil
}
func ImportCardsUpdateCardState(c context.Context, tx *gorm.DB, entity *entity.WSImportCardsEntity) error {
	// 카드 상태 업데이트
	// room_id, card_id, state로 찾고 카드 업데이트할 때 트랜잭션 처리해줘
	for _, card := range entity.Cards {
		err := tx.Model(&mysql.Cards{}).Where("room_id = ? and name = ? and color = ? and state = ?", card.RoomID, card.Name, card.Color, "none").Updates(&mysql.Cards{State: card.State, UserID: card.UserID}).Error
		if err != nil {
			return fmt.Errorf("카드 상태 업데이트 실패 %v", err.Error())
		}
	}
	return nil
}

func ImportCardsUpdateRoomUserCardCount(c context.Context, tx *gorm.DB, entity *entity.WSImportCardsEntity) error {
	// 유저id로 room_users에서 찾아서 card_count를 더한 후 업데이트 한다.
	for _, card := range entity.Cards {
		err := tx.Model(&mysql.RoomUsers{}).Where("room_id = ? AND user_id = ?", card.RoomID, card.UserID).Update("owned_card_count", gorm.Expr("owned_card_count + 1")).Error
		if err != nil {
			return fmt.Errorf("방 유저 카드 카운트 업데이트 실패 %v", err.Error())
		}
	}
	return nil
}

func ImportCardsFindAllCard(c context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.Cards, error) {
	cards := make([]*mysql.Cards, 0)
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and user_id = ?", roomID, userID).Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다. %v", err.Error())
	}
	return cards, nil
}
