package repository

import (
	"context"
	"fmt"
	"main/features/sequence/model/entity"
	_errors "main/features/sequence/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func UseCardFindOneCardInfo(ctx context.Context, tx *gorm.DB, cardID int) (*mysql.SequenceCards, *entity.ErrorInfo) {
	cardInfo := &mysql.SequenceCards{}
	err := tx.Where("id = ?", cardID).First(cardInfo).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 정보 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return cardInfo, nil
}

func UseCardFindOneKingInfo(ctx context.Context, tx *gorm.DB, roomID uint) (int, *entity.ErrorInfo) {
	kingInfo := &mysql.SequenceGameRoomSettings{}
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
func UseCardUpdateKing(ctx context.Context, tx *gorm.DB, roomID uint, kingIndex int) *entity.ErrorInfo {
	// Update both king_index and remaining_slime_count
	err := tx.Model(&mysql.SequenceGameRoomSettings{}).
		Where("room_id = ?", roomID).
		Updates(map[string]interface{}{
			"king_index":            kingIndex,
			"remaining_slime_count": gorm.Expr("remaining_slime_count - ?", 1),
			"current_round":         gorm.Expr("current_round + 1"),
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

func UseCardUpdateUserSlime(ctx context.Context, tx *gorm.DB, roomID, userID uint, nextKingIndex int) *entity.ErrorInfo {
	// Update the slime position for the user
	err := tx.Model(&mysql.SequenceRoomMaps{}).
		Where("room_id = ? AND map_id = ?", roomID, nextKingIndex).
		Updates(map[string]interface{}{
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

func UseCardUpdateCardState(ctx context.Context, tx *gorm.DB, roomID, userID uint, cardID int) *entity.ErrorInfo {
	err := tx.Model(&mysql.SequenceRoomCards{}).
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
