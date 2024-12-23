package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func CloseFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("UserItems").Preload("Room").Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room_users 조회 실패: %v", err),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

func CloseFindOneUser(ctx context.Context, tx *gorm.DB, uID uint) (mysql.Users, *entity.ErrorInfo) {
	var user mysql.Users
	result := tx.WithContext(ctx).Where("id = ?", uID).First(&user)
	if result.Error != nil {
		return mysql.Users{}, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 정보를 찾을 수 없습니다: %v", result.Error),
			Type: _errors.ErrUserNotFound,
		}
	}
	return user, nil
}

func CloseChangeRoomOnwer(ctx context.Context, tx *gorm.DB, RoomID uint, ownerID uint) *entity.ErrorInfo {
	var room mysql.Rooms
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Update("owner_id", ownerID)
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방장 변경 실패: %v", result.Error),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func CloseFindOneRoomUser(ctx context.Context, tx *gorm.DB, RoomID uint) (mysql.RoomUsers, *entity.ErrorInfo) {
	var roomUser mysql.RoomUsers
	result := tx.WithContext(ctx).Where("room_id = ?", RoomID).First(&roomUser)
	if result.Error != nil {
		return mysql.RoomUsers{}, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 정보를 찾을 수 없습니다: %v", result.Error),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUser, nil
}

// 방 삭제
func CloseFindOneAndDeleteRoom(ctx context.Context, tx *gorm.DB, RoomID uint) *entity.ErrorInfo {
	var room mysql.Rooms
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Delete(&room)
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 정보를 삭제할 수 없습니다: %v", result.Error),
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

// 방 유저 삭제
func CloseFindOneAndDeleteRoomUser(ctx context.Context, tx *gorm.DB, uID uint, RoomsID uint) *entity.ErrorInfo {
	var roomUser mysql.RoomUsers
	result := tx.WithContext(ctx).Model(&roomUser).Where("user_id = ? and room_id = ?", uID, RoomsID).Delete(&mysql.RoomUsers{})
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 정보를 삭제할 수 없습니다: %v", result.Error),
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

// 방 정보 업데이트
func CloseFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) (mysql.Rooms, *entity.ErrorInfo) {
	var room mysql.Rooms
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", RoomID).First(&room)
	if result.Error != nil {
		return mysql.Rooms{}, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 정보를 찾을 수 없습니다: %v", result.Error),
			Type: _errors.ErrRoomNotFound,
		}
	}
	room.CurrentCount--
	result = tx.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Updates(room)
	if result.Error != nil {
		return mysql.Rooms{}, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 인원을 업데이트할 수 없습니다: %v", result.Error),
			Type: _errors.ErrUpdateFailed,
		}
	}

	return room, nil
}

// 유저 상태 업데이트
func CloseFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint) *entity.ErrorInfo {
	user := mysql.Users{
		State:  "wait",
		RoomID: 1,
	}
	result := tx.WithContext(ctx).Model(&user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 정보 업데이트 실패: %v", result.Error),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
