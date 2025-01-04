package repository

import (
	"context"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ImportCardsFindAllRoomUsers retrieves all room users with necessary preloads
func ImportCardsFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Table("frog_room_users").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ?", roomID).
		Preload("User").
		Preload("Room").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).
		Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "room_users 조회 실패",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

// ImportCardsFindOneDora retrieves the "dora" card for a specific room
func ImportCardsFindOneDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, *entity.ErrorInfo) {
	var dora mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("state = ?", "dora").
		First(&dora).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "도라 카드를 찾을 수 없습니다.",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return &dora, nil
}

// ImportCardsUpdateCardState updates the state of multiple cards
func ImportCardsUpdateCardState(ctx context.Context, tx *gorm.DB, e *entity.WSImportCardsEntity) *entity.ErrorInfo {
	for _, card := range e.Cards {
		err := tx.Model(&mysql.FrogUserCards{}).
			Where("room_id = ?", card.RoomID).
			Where("card_id = ?", card.CardID).
			Where("state = ?", "none").
			Updates(&mysql.FrogUserCards{
				State:  "owned",
				UserID: card.UserID,
			}).Error
		if err != nil {
			return &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal, // 500
				Msg:  "카드 상태 업데이트 실패",
				Type: _errors.ErrUpdateFailed,
			}
		}
	}
	return nil
}

// ImportCardsUpdateRoomUserCardCount updates the card count for room users
func ImportCardsUpdateRoomUserCardCount(ctx context.Context, tx *gorm.DB, e *entity.WSImportCardsEntity) *entity.ErrorInfo {
	for _, card := range e.Cards {
		err := tx.Model(&mysql.FrogRoomUsers{}).
			Where("room_id = ?", card.RoomID).
			Where("user_id = ?", card.UserID).
			Update("owned_card_count", gorm.Expr("owned_card_count + 1")).Error
		if err != nil {
			return &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal, // 500
				Msg:  "방 유저 카드 카운트 업데이트 실패",
				Type: _errors.ErrUpdateFailed,
			}
		}
	}
	return nil
}

// ImportCardsFindAllCard retrieves all cards for a specific user in a room
func ImportCardsFindAllCard(ctx context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.FrogUserCards, *entity.ErrorInfo) {
	var cards []*mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("user_id = ?", userID).
		Find(&cards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "카드를 찾을 수 없습니다.",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return cards, nil
}
