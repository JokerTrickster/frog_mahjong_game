package repository

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func HeroFindOneCardInfo(ctx context.Context, tx *gorm.DB, roomID uint, cardID int) (*mysql.SlimeWarCards, *entity.ErrorInfo) {
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

func HeroFindOneKingInfo(ctx context.Context, tx *gorm.DB, roomID uint) (int, *entity.ErrorInfo) {
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

func HeroUpdateKing(ctx context.Context, tx *gorm.DB, roomID uint, kingIndex int) *entity.ErrorInfo {
	// Update both king_index and remaining_slime_count
	err := tx.Model(&mysql.SlimeWarGameRoomSettings{}).
		Where("room_id = ?", roomID).
		Updates(map[string]interface{}{
			"king_index":    kingIndex,
			"current_round": gorm.Expr("current_round + 1"),
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

func HeroUpdateUserSlime(ctx context.Context, tx *gorm.DB, roomID uint, uID uint, kingIndex int) *entity.ErrorInfo {
	// Update the slime position for the user
	err := tx.Model(&mysql.SlimeWarRoomMaps{}).
		Where("room_id = ? AND map_id = ?", roomID, kingIndex).
		Updates(map[string]interface{}{
			"user_id": uID,
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

func HeroUpdateCardState(ctx context.Context, tx *gorm.DB, roomID uint, uID uint, cardID int) *entity.ErrorInfo {
	err := tx.Model(&mysql.SlimeWarRoomCards{}).
		Where("room_id = ? AND user_id = ? AND card_id = ?", roomID, uID, cardID).
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

func HeroUpdateUserHeroCardDecrease(ctx context.Context, tx *gorm.DB, roomID uint, uID uint) *entity.ErrorInfo {
	// mysql.SlimeWarUsers
	err := tx.Model(&mysql.SlimeWarUsers{}).
		Where("room_id = ? AND user_id = ?", roomID, uID).
		Update("hero_count", gorm.Expr("hero_count - 1")).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 히어로 카드 감소 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
