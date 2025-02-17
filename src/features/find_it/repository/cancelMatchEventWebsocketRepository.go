package repository

import (
	"context"
	"fmt"
	"main/features/find_it/model/entity"
	_errors "main/features/find_it/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func CancelMatchDeleteRoomUser(ctx context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&mysql.GameRoomUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("roomUser 삭제 실패: %v", err.Error()),
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

func CancelMatchDeleteRoomSetting(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).Delete(&mysql.FindItRoomSettings{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("roomSetting 삭제 실패: %v", err.Error()),
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

func CancelMatchDeleteRoom(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Where("id = ?", roomID).Delete(&mysql.GameRooms{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room 삭제 실패: %v", err.Error()),
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}
