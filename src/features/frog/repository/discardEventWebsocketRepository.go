package repository

import (
	"context"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DiscardCardsFindAllRoomUsers retrieves all room users with necessary preloads
func DiscardCardsFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
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

// DiscardCardsFindOneDora retrieves the "dora" card for the specified room
func DiscardCardsFindOneDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, *entity.ErrorInfo) {
	var dora mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("state = ?", "dora").
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

// DiscardCardsUpdateCardState updates the state of a card to "discard"
func DiscardCardsUpdateCardState(ctx context.Context, tx *gorm.DB, e *entity.WSDiscardCardsEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", e.RoomID).
		Where("card_id = ?", e.CardID).
		Where("state = ?", "owned").
		Updates(map[string]interface{}{
			"state":   "discard",
			"user_id": int(e.UserID),
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

// DiscardCardsUpdateRoomUserCardCount updates the owned card count for a room user
func DiscardCardsUpdateRoomUserCardCount(ctx context.Context, tx *gorm.DB, e *entity.WSDiscardCardsEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("room_id = ?", e.RoomID).
		Where("user_id = ?", e.UserID).
		Update("owned_card_count", gorm.Expr("owned_card_count - 1")).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "방 유저 카드 카운트 업데이트 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// DiscardCardsFindAllCard retrieves all cards owned by a user in a room
func DiscardCardsFindAllCard(ctx context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.FrogUserCards, *entity.ErrorInfo) {
	var cards []*mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("user_id = ?", userID).
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
