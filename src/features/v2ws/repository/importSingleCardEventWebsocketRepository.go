package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	"main/utils/db/mysql"

	_errors "main/features/v2ws/model/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ImportSingleCardFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("User").
		Preload("UserItems").
		Preload("Room").
		Preload("RoomMission").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).
		Preload("UserMissions", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID)
		}).
		Where("room_id = ?", roomID).
		Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room_users 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

func ImportSingleCardUpdateAllCardState(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND state = ?", roomID, "picked").
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

// 카드정보가 존재하면 상태를 업데이트하고 없으면 카드를 생성한다.
func ImportSingleCardCreateCard(ctx context.Context, tx *gorm.DB, userBirdCardDTO *mysql.UserBirdCards) *entity.ErrorInfo {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND card_id = ?", userBirdCardDTO.RoomID, userBirdCardDTO.CardID).
		Updates(userBirdCardDTO).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 조회 및 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func ImportSingleCardUpdateRoomUserCardCount(ctx context.Context, tx *gorm.DB, e *entity.WSImportSingleCardEntity) *entity.ErrorInfo {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.RoomUsers{}).
		Where("room_id = ? AND user_id = ?", e.RoomID, e.UserID).
		Update("owned_card_count", gorm.Expr("owned_card_count + 1")).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 카드 카운트 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func ImportSingleCardFindAllCard(ctx context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.Cards, *entity.ErrorInfo) {
	cards := make([]*mysql.Cards, 0)
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.Cards{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Find(&cards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return cards, nil
}

func ImportSingleCardUpdateOpenCards(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	var count int64
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND state = ?", roomID, "opened").
		Count(&count).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("오픈 카드 카운트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}

	if count != 3 {
		openCardCount := 3 - count
		var cardIDs []int
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&mysql.UserBirdCards{}).
			Where("room_id = ? AND state = ?", roomID, "none").
			Order("RAND()").
			Limit(int(openCardCount)).
			Pluck("card_id", &cardIDs).Error
		if err != nil {
			return &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("카드 조회 실패: %v", err.Error()),
				Type: _errors.ErrUpdateFailed,
			}
		}

		if len(cardIDs) > 0 {
			err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Model(&mysql.UserBirdCards{}).
				Where("room_id = ? AND card_id IN ?", roomID, cardIDs).
				Update("state", "opened").Error
			if err != nil {
				return &entity.ErrorInfo{
					Code: _errors.ErrCodeInternal,
					Msg:  fmt.Sprintf("카드 상태 업데이트 실패: %v", err.Error()),
					Type: _errors.ErrUpdateFailed,
				}
			}
		}
	}

	return nil
}

func ImportSingleCardOwnerCardCount(ctx context.Context, roomID uint, userID uint) (int, *entity.ErrorInfo) {
	var roomUsers mysql.RoomUsers
	err := mysql.GormMysqlDB.Model(&mysql.RoomUsers{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Find(&roomUsers).Error
	if err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 카운트 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}

	if roomUsers.OwnedCardCount > 3 {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeBadRequest,
			Msg:  "카드는 3장을 초과할 수 없습니다.",
			Type: _errors.ErrInvalidRequest,
		}
	}

	return roomUsers.OwnedCardCount, nil
}

func ImportSingleCardFindOneCard(ctx context.Context, roomID uint, cardID uint) *entity.ErrorInfo {
	var card mysql.UserBirdCards
	result := mysql.GormMysqlDB.Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND card_id = ? AND (state = ? OR state = ?)", roomID, cardID, "opened", "none").
		Find(&card)
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 조회 실패: %v", result.Error.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	if result.RowsAffected == 0 {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeBadRequest,
			Msg:  "이미 선택된 카드입니다.",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return nil
}
