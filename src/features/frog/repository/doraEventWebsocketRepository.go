package repository

import (
	"context"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DoraFindAllRoomUsers retrieves all room users with necessary preloads
func DoraFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
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

// DoraCheckFirstPlayer checks if the user is the first player in the room
func DoraCheckFirstPlayer(ctx context.Context, tx *gorm.DB, userID, roomID uint) *entity.ErrorInfo {
	var roomUser mysql.FrogRoomUsers
	err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("user_id = ?", userID).
		Where("room_id = ?", roomID).
		Where("turn_number = ?", 1).
		First(&roomUser).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeForbidden, // 403
			Msg:  "첫 번째 플레이어가 아닙니다",
			Type: _errors.ErrUnauthorizedAction,
		}
	}
	return nil
}

// DoraUpdateDoraCard updates the state of the specified card to "dora"
func DoraUpdateDoraCard(ctx context.Context, tx *gorm.DB, doraEntity *entity.WSDoraEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", doraEntity.RoomID).
		Where("card_id = ?", doraEntity.CardID).
		Update("state", "dora").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "도라 카드 업데이트 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
