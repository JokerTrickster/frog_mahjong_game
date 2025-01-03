package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ImportSingleCardFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
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
func ImportSingleCardFindOneDora(c context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, error) {
	dora := mysql.FrogUserCards{}
	err := tx.Model(&mysql.FrogUserCards{}).Where("room_id = ? and state = ?", roomID, "dora").First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다. %v", err.Error())
	}
	return &dora, nil
}
func ImportSingleCardUpdateCardState(c context.Context, tx *gorm.DB, entity *entity.WSImportSingleCardEntity) error {
	// 카드 상태 업데이트
	// room_id, card_id, state로 찾고 카드 업데이트할 때 트랜잭션 처리해줘
	card := entity.Cards
	err := tx.Model(&mysql.FrogUserCards{}).Where("room_id = ? and card_id = ? and state = ?", card.RoomID, card.CardID, "none").Updates(&mysql.FrogUserCards{State: "owned", UserID: card.UserID}).Error
	if err != nil {
		return fmt.Errorf("카드 상태 업데이트 실패 %v", err.Error())
	}
	return nil
}

func ImportSingleCardUpdateRoomUserCardCount(c context.Context, tx *gorm.DB, entity *entity.WSImportSingleCardEntity) error {
	// 유저id로 room_users에서 찾아서 card_count를 더한 후 업데이트 한다.
	card := entity.Cards
	err := tx.Model(&mysql.FrogRoomUsers{}).Where("room_id = ? AND user_id = ?", card.RoomID, card.UserID).Update("owned_card_count", gorm.Expr("owned_card_count + 1")).Error
	if err != nil {
		return fmt.Errorf("방 유저 카드 카운트 업데이트 실패 %v", err.Error())
	}
	return nil
}

func ImportSingleCardFindAllCard(c context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.FrogUserCards, error) {
	cards := make([]*mysql.FrogUserCards, 0)
	err := tx.Model(&mysql.FrogUserCards{}).Where("room_id = ? and user_id = ?", roomID, userID).Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다. %v", err.Error())
	}
	return cards, nil
}
