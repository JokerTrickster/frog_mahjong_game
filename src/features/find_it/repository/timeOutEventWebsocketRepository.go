package repository

import (
	"context"
	"fmt"
	"main/features/find_it/model/entity"
	_errors "main/features/find_it/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func TimeOutFindCorrectPosition(ctx context.Context, tx *gorm.DB, roomID uint, round, imageID int) ([]*mysql.FindItUserCorrectPositions, *entity.ErrorInfo) {
	var userCorrectPositionDTOList []*mysql.FindItUserCorrectPositions
	err := tx.WithContext(ctx).
		Where("room_id = ?", roomID).
		Where("round = ?", round).
		Where("image_id = ?", imageID).
		Find(&userCorrectPositionDTOList).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("userCorrectPositionDTOList 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return userCorrectPositionDTOList, nil
}

func TimeOutLifeDecrease(ctx context.Context, tx *gorm.DB, roomID uint, diffLife int) *entity.ErrorInfo {
	roomSetting := mysql.FindItRoomSettings{}
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).First(&roomSetting).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	roomSetting.Lifes -= diffLife
	err = tx.WithContext(ctx).Save(&roomSetting).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room 저장 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
