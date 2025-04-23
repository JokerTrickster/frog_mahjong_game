package repository

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
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

func TimeOutFindImageCorrectPosition(ctx context.Context, roomID, round, imageID int) ([]*mysql.FindItImageCorrectPositions, *entity.ErrorInfo) {
	var remainingPositions []*mysql.FindItImageCorrectPositions
	var correctPositionIDs []int

	// ✅ 이미 맞춘 correct_position_id 조회
	err := mysql.GormMysqlDB.WithContext(ctx).
		Table("slime_war_user_correct_positions").
		Select("correct_position_id").
		Where("room_id = ? AND round = ? AND image_id = ?", roomID, round, imageID).
		Pluck("correct_position_id", &correctPositionIDs).Error

	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 정답 위치 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}

	// ✅ 전체 이미지 정답 중에서 맞춘 정답을 제외하고 가져오기
	query := mysql.GormMysqlDB.WithContext(ctx).
		Where("image_id = ?", imageID)

	if len(correctPositionIDs) > 0 {
		query = query.Where("id NOT IN (?)", correctPositionIDs)
	}

	err = query.Find(&remainingPositions).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("남은 정답 위치 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}

	return remainingPositions, nil
}

func TimeOutUpdateTurn(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	if err := tx.Model(&mysql.SlimeWarGameRoomSettings{}).
		Where("room_id = ?", roomID).
		Update("current_count", gorm.Expr("current_count + ?", 1)).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("current_count 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}

	return nil
}
