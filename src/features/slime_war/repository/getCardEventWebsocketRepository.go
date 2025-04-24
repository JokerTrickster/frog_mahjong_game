package repository

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func GetCardCountDummyCard(ctx context.Context, tx *gorm.DB, roomID uint) (int, *entity.ErrorInfo) {
	var count int64
	err := tx.Model(&mysql.SlimeWarRoomCards{}).
		Where("room_id = ? AND state = ?", roomID, "none").Count(&count).Error
	if err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("버린 카드 수 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return int(count), nil
}
func GetCardUpdateDummyCard(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	err := tx.Model(&mysql.SlimeWarRoomCards{}).
		Where("room_id = ? AND state = ?", roomID, "discard").
		Update("state", "none").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("버린 카드 상태 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// 랜덤으로 카드 정보를 하나 가져온다. (state == "none" 인 것 중에)
func GetCardFindOneCardInfo(ctx context.Context, tx *gorm.DB, roomID uint, uID uint) (*mysql.SlimeWarRoomCards, *entity.ErrorInfo) {
	cardInfo := &mysql.SlimeWarRoomCards{}
	err := tx.Where("room_id = ? AND user_id = ? AND state = ?", roomID, uID, "none").Order("RAND()").First(cardInfo).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 정보 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return cardInfo, nil
}

// 카드 상태를 업데이트 한다.
func GetCardUpdateCardState(ctx context.Context, tx *gorm.DB, roomID uint, uID uint, cardID int) *entity.ErrorInfo {
	err := tx.Model(&mysql.SlimeWarRoomCards{}).
		Where("room_id = ? AND user_id = ? AND card_id = ?", roomID, uID, cardID).
		Update("state", "owned").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 상태 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// 라운드 수를 업데이트 한다. 남은 카드 수를 감소시킨다.
func GetCardUpdateRoomSetting(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	err := tx.Model(&mysql.SlimeWarGameRoomSettings{}).
		Where("room_id = ?", roomID).
		Update("current_round", gorm.Expr("current_round + 1")).
		Update("remaining_card_count", gorm.Expr("remaining_card_count - 1")).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("라운드 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
