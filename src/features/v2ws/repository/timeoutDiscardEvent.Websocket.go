package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func TimeOutDiscardCardsFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	err := tx.Preload("User").
		Preload("Room").
		Preload("RoomMission").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).
		Preload("UserMissions", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID)
		}).
		Where("room_id = ?", roomID).
		Find(&roomUsers).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room_users 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

func TimeOutDiscardCardsFindOneDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.Cards, *entity.ErrorInfo) {
	dora := mysql.Cards{}
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? AND state = ?", roomID, "dora").First(&dora).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  fmt.Sprintf("도라 카드를 찾을 수 없습니다. %v", err.Error()),
			Type: _errors.ErrNotFoundCard,
		}
	}
	return &dora, nil
}

func TimeOutDiscardCardsUpdateCardState(ctx context.Context, tx *gorm.DB, e *entity.WSTimeOutDiscardCardsEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.Cards{}).
		Where("room_id = ? AND card_id = ? AND state = ?", e.RoomID, e.CardID, "none").
		Updates(&mysql.Cards{State: "discard", UserID: int(e.UserID)}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 버리기 상태 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func TimeOutDiscardCardsFindAllCard(ctx context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.Cards, *entity.ErrorInfo) {
	cards := make([]*mysql.Cards, 0)
	err := tx.Model(&mysql.Cards{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Order("updated_at ASC").
		Find(&cards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드를 찾을 수 없습니다. %v", err.Error()),
			Type: _errors.ErrNotFoundCard,
		}
	}
	return cards, nil
}
