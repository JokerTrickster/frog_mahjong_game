package repository

import (
	"context"
	"main/features/frog/model/entity"
	"main/utils/db/mysql"

	_errors "main/features/frog/model/errors"

	"gorm.io/gorm"
)

// CloseFindAllRoomUsers retrieves all room users with necessary preloads
func CloseFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").
		Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  "room_users 조회 실패",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

// CloseFindOneUser retrieves user information by ID
func CloseFindOneUser(ctx context.Context, tx *gorm.DB, uID uint) (mysql.GameUsers, *entity.ErrorInfo) {
	var user mysql.GameUsers
	if err := tx.WithContext(ctx).Where("id = ?", uID).First(&user).Error; err != nil {
		return mysql.GameUsers{}, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  "유저 정보를 찾을 수 없습니다",
			Type: _errors.ErrUserNotFound,
		}
	}
	return user, nil
}

// CloseChangeRoomOwner updates the owner of a room
func CloseChangeRoomOnwer(ctx context.Context, tx *gorm.DB, RoomID uint, ownerID uint) *entity.ErrorInfo {
	if err := tx.WithContext(ctx).
		Model(&mysql.GameRooms{}).
		Where("id = ?", RoomID).
		Update("owner_id", ownerID).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "방장 변경 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// CloseFindOneRoomUser retrieves a single room user by room ID
func CloseFindOneRoomUser(ctx context.Context, tx *gorm.DB, RoomID uint) (mysql.FrogRoomUsers, *entity.ErrorInfo) {
	var roomUser mysql.FrogRoomUsers
	if err := tx.WithContext(ctx).Where("room_id = ?", RoomID).First(&roomUser).Error; err != nil {
		return mysql.FrogRoomUsers{}, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  "방 유저 정보를 찾을 수 없습니다",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUser, nil
}

// CloseFindOneAndDeleteRoom deletes a room by ID
func CloseFindOneAndDeleteRoom(ctx context.Context, tx *gorm.DB, RoomID uint) *entity.ErrorInfo {
	if err := tx.WithContext(ctx).
		Where("id = ?", RoomID).
		Delete(&mysql.GameRooms{}).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "방 정보를 삭제할 수 없습니다",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

// CloseFindOneAndDeleteRoomUser deletes a user from a room
func CloseFindOneAndDeleteRoomUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) *entity.ErrorInfo {
	if err := tx.WithContext(ctx).
		Where("user_id = ? and room_id = ?", uID, RoomID).
		Delete(&mysql.FrogRoomUsers{}).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "방 유저 정보를 삭제할 수 없습니다",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

// CloseFindOneAndUpdateRoom decreases room user count by 1
func CloseFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) (mysql.GameRooms, *entity.ErrorInfo) {
	var room mysql.GameRooms
	if err := tx.WithContext(ctx).Where("id = ?", RoomID).First(&room).Error; err != nil {
		return mysql.GameRooms{}, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  "방 정보를 찾을 수 없습니다",
			Type: _errors.ErrRoomNotFound,
		}
	}
	room.CurrentCount--
	if err := tx.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Updates(room).Error; err != nil {
		return mysql.GameRooms{}, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "방 인원을 업데이트할 수 없습니다",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return room, nil
}

// CloseFindOneAndUpdateUser updates user information when they leave a room
func CloseFindOneAndUpdateUser(ctx context.Context, uID uint) *entity.ErrorInfo {
	user := &mysql.GameUsers{
		State:  "wait",
		RoomID: 1, // Set default RoomID
	}
	if err := mysql.GormMysqlDB.WithContext(ctx).Model(&user).Where("id = ?", uID).Updates(&user).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "유저 정보 업데이트 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
