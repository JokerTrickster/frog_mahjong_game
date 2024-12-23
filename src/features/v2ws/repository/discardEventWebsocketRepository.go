package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func DiscardCardsFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("User").
		Preload("UserItems").
		Preload("Room").
		Preload("RoomMission").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).
		Preload("UserMissions", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID)
		}).
		Where("room_id = ?", roomID).
		Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room_users 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

func DiscardCardUpdateAllCardState(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND state = ?", roomID, "picked").
		Update("state", "discard").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 상태 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func DiscardCardsUpdateCardState(ctx context.Context, tx *gorm.DB, e *entity.WSDiscardCardsEntity) *entity.ErrorInfo {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND card_id = ? AND state = ?", e.RoomID, e.CardID, "owned").
		Updates(&mysql.Cards{State: "picked", UserID: int(e.UserID)}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 상태 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func DiscardCardsUpdateRoomUserCardCount(ctx context.Context, tx *gorm.DB, e *entity.WSDiscardCardsEntity) *entity.ErrorInfo {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.RoomUsers{}).
		Where("room_id = ? AND user_id = ?", e.RoomID, e.UserID).
		Update("owned_card_count", gorm.Expr("owned_card_count - 1")).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 카드 카운트 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func DiscardCardsOwnerCardCount(ctx context.Context, roomID uint, userID uint) (int, *entity.ErrorInfo) {
	var roomUsers mysql.RoomUsers
	err := mysql.GormMysqlDB.Model(&mysql.RoomUsers{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Find(&roomUsers).Error
	if err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 카운트 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers.OwnedCardCount, nil
}
