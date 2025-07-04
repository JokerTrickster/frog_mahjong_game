package repository

import (
	"context"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AbnormalFindAllRoomUsers retrieves all room users with necessary preloads
func AbnormalFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Table("frog_room_users").Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ?", roomID).
		Preload("User").
		Preload("Room").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  "room_users 조회 실패",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

// AbnormalDeleteAllCards deletes all cards for a given room
func AbnormalDeleteAllCards(ctx context.Context, tx *gorm.DB, abnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	if err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", abnormalEntity.RoomID).
		Delete(&mysql.FrogUserCards{}).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "카드 삭제 실패",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

// AbnormalDeleteRoom deletes a room by room ID
func AbnormalDeleteRoom(ctx context.Context, tx *gorm.DB, abnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	if err := tx.Model(&mysql.GameRooms{}).
		Where("id = ?", abnormalEntity.RoomID).
		Delete(&mysql.GameRooms{}).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "방 삭제 실패",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

// AbnormalUpdateUsers updates the state of users in a room to "wait"
func AbnormalUpdateUsers(ctx context.Context, tx *gorm.DB, abnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	if err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("room_id = ?", abnormalEntity.RoomID).
		Where("user_id = ?", abnormalEntity.AbnormalUserID).
		Update("player_state", "disconnected").Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "유저 상태 변경 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func AbnormalDeleteRoomUsers(ctx context.Context, tx *gorm.DB, abnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	if err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("room_id = ?", abnormalEntity.RoomID).
		Delete(&mysql.FrogRoomUsers{}).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "방 유저 정보 삭제 실패",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

func AbnormalDeleteGameRoomSetting(ctx context.Context, tx *gorm.DB, abnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	if err := tx.Model(&mysql.FrogGameRoomSettings{}).
		Where("room_id = ?", abnormalEntity.RoomID).
		Delete(&mysql.FrogGameRoomSettings{}).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "방 게임 셋팅 삭제 실패",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}
