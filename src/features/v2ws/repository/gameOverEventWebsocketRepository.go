package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func GameOverFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").
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

// 유저 상태 변경 (play -> wait)
func GameOverUpdateRoomUsers(ctx context.Context, tx *gorm.DB, GameOverEntity *entity.WSGameOverEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.RoomUsers{}).
		Where("room_id = ?", GameOverEntity.RoomID).
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
