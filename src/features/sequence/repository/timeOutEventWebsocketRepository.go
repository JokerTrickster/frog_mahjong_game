package repository

import (
	"context"
	"fmt"
	"main/features/sequence/model/entity"
	_errors "main/features/sequence/model/errors"
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
		Table("sequence_user_correct_positions").
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
	if err := tx.Model(&mysql.SequenceGameRoomSettings{}).
		Where("room_id = ?", roomID).
		Update("current_round", gorm.Expr("current_round + ?", 1)).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("current_count 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}

	return nil
}
func TimeOutFindUserCards(ctx context.Context, tx *gorm.DB, roomID, userID uint) ([]*mysql.SequenceRoomCards, *entity.ErrorInfo) {
	var userCards []*mysql.SequenceRoomCards
	err := tx.WithContext(ctx).Where("room_id = ? and user_id = ? and state = ?", roomID, userID, "owned").Find(&userCards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("userCards 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return userCards, nil
}

func TimeOutUpdateCardState(ctx context.Context, tx *gorm.DB, roomID, cardID int) *entity.ErrorInfo {
	if err := tx.Model(&mysql.SequenceRoomCards{}).
		Where("room_id = ? and card_id = ? and state = ?", roomID, cardID, "owned").
		Update("state", "used").Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("cardState 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func TimeOutUpdateMapState(ctx context.Context, tx *gorm.DB, roomID, userID, cardID int) *entity.ErrorInfo {
	if err := tx.Model(&mysql.SequenceRoomMaps{}).
		Where("room_id = ? and map_id = ? and user_id = ?", roomID, cardID, 0).
		Update("user_id", userID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("mapState 업데이트 실패: %v", err.Error()),
				Type: _errors.ErrUpdateFailed,
			}
		}
	}
	return nil
}

func TimeOutUpdateDummyCardState(ctx context.Context, tx *gorm.DB, roomID, userID int) *entity.ErrorInfo {
	// 랜덤으로 한 장의 카드만 가져오기
	err := tx.Model(&mysql.SequenceRoomCards{}).
		Where("room_id = ? AND state = ?", roomID, "none").
		Order("RAND()"). // 랜덤 정렬
		Limit(1).        // 한 장만 선택
		Updates(map[string]interface{}{
			"user_id": userID,
			"state":   "owned",
		}).Error

	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("더미 카드 상태 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
