package repository

import (
	"context"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	_errors "main/features/ws/model/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TimeOutDiscardCardsFindAllRoomUsers retrieves all room users with necessary preloads
func TimeOutDiscardCardsFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
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
			Msg:  "방 유저를 찾을 수 없습니다",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

// TimeOutDiscardFindDora retrieves the "dora" card for a specific room
func TimeOutDiscardFindDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, *entity.ErrorInfo) {
	var dora mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ? AND state = ?", roomID, "dora").
		First(&dora).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "도라 카드를 찾을 수 없습니다",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return &dora, nil
}

// TimeOutDiscardUpdateCardState updates the state of a specific card
func TimeOutDiscardUpdateCardState(ctx context.Context, tx *gorm.DB, e *entity.WSTimeOutDiscardCardsEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ? AND card_id = ? AND state = ?", e.RoomID, e.CardID, "none").
		Updates(&mysql.FrogUserCards{
			State:  "discard",
			UserID: int(e.UserID),
		}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "카드 상태 업데이트 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// TimeOutDiscardFindAllCards retrieves all cards for a specific room and user
func TimeOutDiscardFindAllCards(ctx context.Context, tx *gorm.DB, roomID, userID uint) ([]*mysql.FrogUserCards, *entity.ErrorInfo) {
	var cards []*mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Order("updated_at ASC").
		Find(&cards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "카드를 찾을 수 없습니다",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return cards, nil
}
