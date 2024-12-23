package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func CancelMatchFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room_users 조회 실패: %v", err),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

func CancelMatchDeleteOneRoomUser(ctx context.Context, tx *gorm.DB, roomID, uID uint) *entity.ErrorInfo {
	roomUser := mysql.RoomUsers{}
	err := tx.Where("room_id = ? AND user_id = ?", roomID, uID).Delete(&roomUser).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 삭제 실패: %v", err),
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

func CancelMatchFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.Rooms, *entity.ErrorInfo) {
	var room mysql.Rooms
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).First(&room)
	if result.Error != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 정보를 찾을 수 없습니다: %v", result.Error),
			Type: _errors.ErrRoomNotFound,
		}
	}
	room.CurrentCount--
	result = tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).Updates(room)
	if result.Error != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 인원을 업데이트할 수 없습니다: %v", result.Error),
			Type: _errors.ErrUpdateFailed,
		}
	}

	if room.CurrentCount == 0 {
		result = tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).Delete(&room)
		if result.Error != nil {
			return nil, &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("방 정보를 삭제할 수 없습니다: %v", result.Error),
				Type: _errors.ErrDeleteFailed,
			}
		}
	}
	return &room, nil
}

func CancelMatchFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint) *entity.ErrorInfo {
	user := &mysql.Users{
		RoomID: 1,
		State:  "wait",
	}
	err := tx.Model(&user).Where("id = ?", uID).Updates(user).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 정보 업데이트 실패: %v", err),
			Type: _errors.ErrUpdateFailed,
		}
	}

	return nil
}

func CancelMatchFindOneRoomUser(ctx context.Context, tx *gorm.DB, roomID uint) (uint, *entity.ErrorInfo) {
	var roomUser mysql.RoomUsers
	result := tx.WithContext(ctx).Where("room_id = ?", roomID).First(&roomUser)
	if result.Error != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 정보를 찾을 수 없습니다: %v", result.Error),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return uint(roomUser.UserID), nil
}

func CancelMatchUpdateRoomOwner(ctx context.Context, tx *gorm.DB, roomID uint, roomUserID uint) *entity.ErrorInfo {
	var room mysql.Rooms
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).Update("owner_id", roomUserID)
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방장 변경 실패: %v", result.Error),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
