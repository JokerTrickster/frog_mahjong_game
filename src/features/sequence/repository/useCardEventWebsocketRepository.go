package repository

import (
	"context"
	"fmt"
	"main/features/sequence/model/entity"
	_errors "main/features/sequence/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func UseCardUpdateCardState(ctx context.Context, tx *gorm.DB, roomID, userID uint, cardID int) *entity.ErrorInfo {
	err := tx.Model(&mysql.SequenceRoomCards{}).
		Where("room_id = ? AND user_id = ? AND card_id = ? AND state = ?", roomID, userID, cardID, "owned").
		Updates(map[string]interface{}{
			"state": "used",
		}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 상태 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func UseCardUpdateMapState(ctx context.Context, tx *gorm.DB, roomID, userID, mapID int) *entity.ErrorInfo {
	err := tx.Model(&mysql.SequenceRoomMaps{}).
		Where("room_id = ? AND  map_id = ? AND user_id = ?", roomID, mapID, 0).
		Updates(map[string]interface{}{
			"user_id": userID,
		}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("맵 상태 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func UseCardGetDummyCard(ctx context.Context, tx *gorm.DB, roomID, userID uint) *entity.ErrorInfo {
	err := tx.Model(&mysql.SequenceRoomCards{}).
		Where("room_id = ? AND state = ?", roomID, "none").
		Updates(map[string]interface{}{
			"state":   "owned",
			"user_id": userID,
		}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("더미 카드 가져오기 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
