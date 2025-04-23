package repository

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func MoveFindOneCardInfo(ctx context.Context, tx *gorm.DB, roomID uint, cardID int) (*mysql.SlimeWarCards, *entity.ErrorInfo) {
	cardInfo := &mysql.SlimeWarCards{}
	err := tx.Where("id = ?", roomID, cardID).First(cardInfo).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 정보 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return cardInfo, nil
}

func MoveFindOneKingInfo(ctx context.Context, tx *gorm.DB, roomID uint) (int, *entity.ErrorInfo) {
	kingInfo := &mysql.SlimeWarGameRoomSettings{}
	err := tx.Where("room_id = ?", roomID).First(kingInfo).Error
	if err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("왕 정보 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return kingInfo.KingIndex, nil
}

// 왕 정보 업데이트
func MoveUpdateKing(ctx context.Context, tx *gorm.DB, roomID uint, kingIndex int) *entity.ErrorInfo {
	// Update both king_index and remaining_slime_count
	err := tx.Model(&mysql.SlimeWarGameRoomSettings{}).
		Where("room_id = ?", roomID).
		Updates(map[string]interface{}{
			"king_index":            kingIndex,
			"remaining_slime_count": gorm.Expr("remaining_slime_count - ?", 1),
		}).Error

	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("왕 정보 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func MoveUpdateUserSlime(ctx context.Context, tx *gorm.DB, roomID, userID uint, nextKingIndex int) *entity.ErrorInfo {
	// Update the slime position for the user
	err := tx.Model(&mysql.SlimeWarRoomCards{}).
		Where("room_id = ? AND map_id = ?", roomID, nextKingIndex).
		Updates(map[string]interface{}{
			"state":   "owned",
			"user_id": userID,
		}).Error

	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 슬라임 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func MoveUpdateCardState(ctx context.Context, tx *gorm.DB, roomID, userID uint, cardID int) *entity.ErrorInfo {
	err := tx.Model(&mysql.SlimeWarRoomCards{}).
		Where("room_id = ? AND user_id = ? AND card_id = ?", roomID, userID, cardID).
		Updates(map[string]interface{}{
			"state":   "discard",
			"user_id": 0,
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
