package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func RequestWinFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	err := tx.Preload("User").
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

// 카드 정보 체크 (소유하고 있는지 체크)
func RequestWinFindAllCards(ctx context.Context, tx *gorm.DB, requestWinEntity *entity.V2WSRequestWinEntity) ([]*mysql.UserBirdCards, *entity.ErrorInfo) {
	cards := make([]*mysql.UserBirdCards, 0)
	err := tx.Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND user_id = ? AND card_id IN ?", requestWinEntity.RoomID, requestWinEntity.UserID, requestWinEntity.Cards).
		Find(&cards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 조회 실패: %v", err.Error()),
			Type: _errors.ErrNotFoundCard,
		}
	}
	return cards, nil
}

// 유저 상태 변경 (play -> wait)
func RequestWinUpdateRoomUsers(ctx context.Context, tx *gorm.DB, requestWinEntity *entity.V2WSRequestWinEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.RoomUsers{}).
		Where("room_id = ?", requestWinEntity.RoomID).
		Update("player_state", "wait").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 상태 변경 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
