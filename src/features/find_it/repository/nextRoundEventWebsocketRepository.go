package repository

import (
	"context"
	"main/features/find_it/model/entity"
	_errors "main/features/find_it/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NextRoundRoundIncrease(ctx context.Context, tx *gorm.DB, roomID uint, round int) *entity.ErrorInfo {
	// 라운드 증가
	rooomSetting := mysql.FindItRoomSettings{}
	if err := tx.Model(&rooomSetting).Where("id = ?", roomID).Update("round", round).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "라운드 증가 실패",
			Type: _errors.ErrRoundIncreaseFailed,
		}
	}
	return nil
}
