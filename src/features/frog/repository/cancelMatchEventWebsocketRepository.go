package repository

import (
	"context"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CancelMatchFindAllRoomUsers retrieves all room users with necessary preloads
func CancelMatchFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
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

// CancelMatchDeleteOneRoomUser deletes a room user by room ID and user ID
func CancelMatchDeleteOneRoomUser(ctx context.Context, tx *gorm.DB, roomID, uID uint) *entity.ErrorInfo {
	if err := tx.Where("room_id = ? AND user_id = ?", roomID, uID).Delete(&mysql.FrogRoomUsers{}).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "방 유저 삭제 실패",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

// CancelMatchFindOneAndUpdateRoom updates room information and deletes it if empty
func CancelMatchFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.GameRooms, *entity.ErrorInfo) {
	var room mysql.GameRooms
	if err := tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).First(&room).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  "방 정보를 찾을 수 없습니다",
			Type: _errors.ErrRoomNotFound,
		}
	}
	room.CurrentCount--
	if err := tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).Updates(room).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "방 인원을 업데이트할 수 없습니다",
			Type: _errors.ErrUpdateFailed,
		}
	}
	if room.CurrentCount == 0 {
		if err := tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).Delete(&room).Error; err != nil {
			return nil, &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  "방 삭제 실패",
				Type: _errors.ErrDeleteFailed,
			}
		}
	}
	return &room, nil
}

// CancelMatchFindOneAndUpdateUser updates user information to default state
func CancelMatchFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint) *entity.ErrorInfo {
	user := &mysql.GameUsers{
		RoomID: 1,
		State:  "wait",
	}
	if err := tx.Model(user).Where("id = ?", uID).Updates(user).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "유저 정보 업데이트 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// CancelMatchFindOneRoomUser retrieves the first room user for a given room ID
func CancelMatchFindOneRoomUser(ctx context.Context, tx *gorm.DB, roomID uint) (uint, *entity.ErrorInfo) {
	var roomUser mysql.FrogRoomUsers
	if err := tx.WithContext(ctx).Where("room_id = ?", roomID).First(&roomUser).Error; err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  "방 유저 정보를 찾을 수 없습니다",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return uint(roomUser.UserID), nil
}

// CancelMatchUpdateRoomOwner updates the owner of the room
func CancelMatchUpdateRoomOwner(ctx context.Context, tx *gorm.DB, roomID uint, roomUserID uint) *entity.ErrorInfo {
	if err := tx.WithContext(ctx).Model(&mysql.GameRooms{}).Where("id = ?", roomID).Update("owner_id", roomUserID).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "방장 변경 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
