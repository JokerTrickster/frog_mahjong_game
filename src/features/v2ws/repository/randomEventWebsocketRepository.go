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

func RandomFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	err := tx.Preload("User").
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
		Find(&roomUsers).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room_users 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

func RandomUpdateRandomCards(ctx context.Context, tx *gorm.DB, randomEntity *entity.WSRandomEntity) *entity.ErrorInfo {
	var cardIDs []int
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND state IN ?", randomEntity.RoomID, []string{"opened", "none"}).
		Order("RAND()").
		Limit(int(randomEntity.Count)).
		Pluck("card_id", &cardIDs).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 조회 실패: %v", err.Error()),
			Type: _errors.ErrNotFoundCard,
		}
	}

	if len(cardIDs) > 0 {
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&mysql.UserBirdCards{}).
			Where("room_id = ? AND user_id = ? AND card_id IN ?", randomEntity.RoomID, 0, cardIDs).
			Updates(map[string]interface{}{
				"state":   "owned",
				"user_id": randomEntity.UserID,
			}).Error
		if err != nil {
			return &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("카드 상태 업데이트 실패: %v", err.Error()),
				Type: _errors.ErrUpdateFailed,
			}
		}
	}

	return nil
}

func RandomCreateCard(ctx context.Context, tx *gorm.DB, userBirdCardDTO *mysql.UserBirdCards) *entity.ErrorInfo {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND card_id = ?", userBirdCardDTO.RoomID, userBirdCardDTO.CardID).
		Updates(userBirdCardDTO).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 생성 또는 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrCreateFailed,
		}
	}

	return nil
}

func RandomUpdateRoomUserCardCount(ctx context.Context, tx *gorm.DB, e *entity.WSRandomEntity) *entity.ErrorInfo {
	addedCount := fmt.Sprintf("owned_card_count + %d", e.Count)
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.RoomUsers{}).
		Where("room_id = ? AND user_id = ?", e.RoomID, e.UserID).
		Update("owned_card_count", gorm.Expr(addedCount)).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 카드 카운트 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func RandomFindAllCard(ctx context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.Cards, *entity.ErrorInfo) {
	cards := make([]*mysql.Cards, 0)
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.Cards{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Find(&cards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 조회 실패: %v", err.Error()),
			Type: _errors.ErrNotFoundCard,
		}
	}
	return cards, nil
}

func RandomUpdateOpenCards(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	var count int64
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND state = ?", roomID, "opened").
		Count(&count).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("오픈 카드 카운트 조회 실패: %v", err.Error()),
			Type: _errors.ErrCountFailed,
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
				Type: _errors.ErrNotFoundCard,
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

func RandomUpdateAllCardState(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
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
