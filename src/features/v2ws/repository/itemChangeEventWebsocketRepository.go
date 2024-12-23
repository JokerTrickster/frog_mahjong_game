package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/utils/db/mysql"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func ItemChangeFindAllRoomUsers(ctx context.Context, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := mysql.GormMysqlDB.Preload("User").
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

func ItemChangeCheck(ctx context.Context, itemChangeEntity entity.WSItemChangeEntity) *entity.ErrorInfo {
	var userItems mysql.UserItems
	result := mysql.GormMysqlDB.Model(&mysql.UserItems{}).
		Where("room_id = ? AND user_id = ? AND item_id = ?", itemChangeEntity.RoomID, itemChangeEntity.UserID, itemChangeEntity.ItemID).
		Find(&userItems)
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("아이템 조회 실패: %v", result.Error.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	if userItems.RemainingUses == 0 {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeBadRequest,
			Msg:  "아이템 사용 횟수가 없습니다.",
			Type: _errors.ErrItemNotAvailable,
		}
	}
	return nil
}

func ItemChange(ctx context.Context, tx *gorm.DB, itemChangeEntity entity.WSItemChangeEntity) *entity.ErrorInfo {
	roomID := itemChangeEntity.RoomID

	var noneCards []mysql.UserBirdCards
	if err := tx.Where("room_id = ? AND state = ?", roomID, "none").
		Find(&noneCards).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("none 상태 카드 조회 중 에러 발생: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}

	if err := tx.Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND state = ?", roomID, "opened").
		Update("state", "none").Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 상태 변경 (opened -> none) 중 에러 발생: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}

	noneCardCount := len(noneCards)
	if noneCardCount < 3 {
		var noneCardIDs []int
		for _, card := range noneCards {
			noneCardIDs = append(noneCardIDs, card.CardID)
		}

		if err := tx.Model(&mysql.UserBirdCards{}).
			Where("room_id = ? AND id IN ?", itemChangeEntity.RoomID, noneCardIDs).
			Update("state", "opened").Error; err != nil {
			return &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("카드 상태 변경 (none -> opened) 중 에러 발생: %v", err.Error()),
				Type: _errors.ErrUpdateFailed,
			}
		}
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(noneCardCount, func(i, j int) {
		noneCards[i], noneCards[j] = noneCards[j], noneCards[i]
	})

	selectedCards := noneCards[:3]
	var selectedCardIDs []int
	for _, card := range selectedCards {
		selectedCardIDs = append(selectedCardIDs, card.CardID)
	}

	if err := tx.Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND card_id IN ?", itemChangeEntity.RoomID, selectedCardIDs).
		Update("state", "opened").Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 상태 변경 (none -> opened) 중 에러 발생: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}

	return nil
}

func ItemChangeConsumeUserItems(ctx context.Context, tx *gorm.DB, itemChangeEntity entity.WSItemChangeEntity) *entity.ErrorInfo {
	if err := tx.Model(&mysql.UserItems{}).
		Where("room_id = ? AND user_id = ? AND item_id = ?", itemChangeEntity.RoomID, itemChangeEntity.UserID, itemChangeEntity.ItemID).
		Update("remaining_uses", gorm.Expr("remaining_uses - 1")).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("user_item remaining_uses -1 중 에러 발생: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}

	return nil
}
