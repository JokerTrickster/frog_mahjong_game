package repository

import (
	"context"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GameOverFindAllRoomUsers retrieves all room users with necessary preloads
func GameOverFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
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

// GameOverDeleteAllCards deletes all cards associated with a specific room
func GameOverDeleteAllCards(ctx context.Context, tx *gorm.DB, gameOverEntity *entity.WSGameOverEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", gameOverEntity.RoomID).
		Delete(&mysql.FrogUserCards{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "카드 삭제 실패",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

// GameOverUpdateRoomUsers updates the state of all users in a room to "wait"
func GameOverUpdateRoomUsers(ctx context.Context, tx *gorm.DB, gameOverEntity *entity.WSGameOverEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("room_id = ?", gameOverEntity.RoomID).
		Update("player_state", "wait").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "방 유저 상태 변경 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
