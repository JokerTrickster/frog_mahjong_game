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

func SubmitPositionCheck(ctx context.Context, tx *gorm.DB, imageID int, xPosition, yPosition float64) (int, *entity.ErrorInfo) {
	const threshold = 10.0 // ✅ 허용 오차 (10px)

	// find_it_image_correct_positions
	var correctPosition mysql.FindItImageCorrectPositions
	err := tx.WithContext(ctx).
		Where("image_id = ? AND ABS(x_position - ?) <= ? AND ABS(y_position - ?) <= ?",
			imageID, xPosition, threshold, yPosition, threshold).
		First(&correctPosition).Error

	// ✅ 좌표가 없으면 false 반환
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // 좌표 없음
		}
		// 기타 DB 오류 발생
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DB 오류 발생: %v", err.Error()),
			Type: _errors.ErrDBServer,
		}
	}

	// ✅ 좌표가 있으면 true 반환
	return int(correctPosition.ID), nil
}

func SubmitPositionCorrectSave(ctx context.Context, tx *gorm.DB, userCorrectPosition *mysql.FindItUserCorrectPositions) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Create(userCorrectPosition).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("정답 좌표 저장 실패: %v", err.Error()),
			Type: _errors.ErrCreateFailed,
		}
	}
	return nil
}

func SubmitPositionLifeDecrease(ctx context.Context, tx *gorm.DB, roomID uint, uID uint) *entity.ErrorInfo {
	// find_it_game_users
	var roomSettings mysql.FindItRoomSettings
	err := tx.WithContext(ctx).
		Where("room_id = ?", roomID).
		First(&roomSettings).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("게임 유저 정보 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	// 목숨 감소
	err = tx.WithContext(ctx).
		Model(&roomSettings).
		Update("lifes", gorm.Expr("lifes - 1")).
		Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("목숨 감소 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}

	return nil
}
