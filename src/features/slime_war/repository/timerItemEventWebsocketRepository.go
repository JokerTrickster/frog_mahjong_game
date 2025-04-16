package repository

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func TimerItemCheck(ctx context.Context, roomID uint) (*mysql.FindItRoomSettings, *entity.ErrorInfo) {
	var roomSettings mysql.FindItRoomSettings
	err := mysql.GormMysqlDB.WithContext(ctx).
		Where("room_id = ?", roomID).
		First(&roomSettings).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("roomSettings 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return &roomSettings, nil
}

func TimerItemDecrease(ctx context.Context, tx *gorm.DB, roomSettingDTO *mysql.FindItRoomSettings) *entity.ErrorInfo {
	roomSettingDTO.ItemTimerStopCount--
	err := tx.WithContext(ctx).Save(&roomSettingDTO).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("타이머 아이템 1 감소 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
