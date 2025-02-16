package repository

import (
	"context"
	"errors"
	"fmt"
	"main/features/find_it/model/entity"
	_errors "main/features/find_it/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func HintItemCheck(ctx context.Context, roomID uint) (*mysql.FindItRoomSettings, *entity.ErrorInfo) {
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

func HintItemDecrease(ctx context.Context, tx *gorm.DB, roomSettingDTO *mysql.FindItRoomSettings) *entity.ErrorInfo {
	roomSettingDTO.ItemHintCount--
	err := tx.WithContext(ctx).Save(&roomSettingDTO).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("힌트 아이템 1 감소 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func HintItemFindCorrectPosition(ctx context.Context, tx *gorm.DB, roomID uint, round, imageID int) ([]*mysql.FindItUserCorrectPositions, *entity.ErrorInfo) {
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

func HintItemFindOneCorrectPosition(ctx context.Context, tx *gorm.DB, imageID uint, correctPositionIDList []uint) (*mysql.FindItImageCorrectPositions, *entity.ErrorInfo) {
	var correctPosition mysql.FindItImageCorrectPositions

	// 만약 제외할 리스트가 없으면 그냥 랜덤으로 하나 선택
	query := tx.WithContext(ctx).Where("image_id = ?", imageID)
	if len(correctPositionIDList) > 0 {
		query = query.Where("id NOT IN (?)", correctPositionIDList)
	}

	// 랜덤으로 하나 가져오기
	err := query.Order("RAND()").Limit(1).First(&correctPosition).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 정답 좌표 없음
		}
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DB 오류 발생: %v", err.Error()),
			Type: _errors.ErrDBServer,
		}
	}

	return &correctPosition, nil
}
