package repository

import (
	"context"
	"fmt"
	"main/features/sequence/model/entity"
	_errors "main/features/sequence/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func GameOverFindGameRoomUser(ctx context.Context, tx *gorm.DB, roomID uint) ([]mysql.GameRoomUsers, *entity.ErrorInfo) {
	roomUsers := []mysql.GameRoomUsers{}
	if err := tx.Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("GameOverFindGameRoomUser: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return roomUsers, nil
}

func GameOverSaveGameResult(ctx context.Context, tx *gorm.DB, gameResult mysql.GameResults) *entity.ErrorInfo {
	if err := tx.Create(&gameResult).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("GameOverSaveGameResult: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}
