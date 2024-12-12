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

func ItemChangeFindAllRoomUsers(ctx context.Context, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := mysql.GormMysqlDB.Preload("User").Preload("UserItems").Preload("Room").Preload("RoomMission").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Preload("UserMissions", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID)
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}

func ItemChangeCheck(ctx context.Context, itemChangeEntity entity.WSItemChangeEntity) *entity.ErrorInfo {
	var userItems mysql.UserItems

	result := mysql.GormMysqlDB.Model(&mysql.UserItems{}).Where("room_id = ? and user_id = ? and item_id = ?", itemChangeEntity.RoomID, itemChangeEntity.UserID, itemChangeEntity.ItemID).Find(&userItems)
	if result.Error != nil {
		return &entity.ErrorInfo{Code: _errors.ErrCodeInternal, Msg: result.Error.Error(), Type: _errors.ErrInternalServer}
	}
	if userItems.RemainingUses == 0 {
		return &entity.ErrorInfo{Code: _errors.ErrCodeBadRequest, Msg: "아이템 사용 횟수가 없습니다.", Type: _errors.ErrItemNotAvailable}
	}
	return nil
}

func ItemChange(ctx context.Context, tx *gorm.DB, itemChangeEntity entity.WSItemChangeEntity) error {
	roomID := itemChangeEntity.RoomID
	// 2. room_id로 조회하여 state가 "none"인 카드들을 가져오기
	var noneCards []mysql.UserBirdCards
	if err := tx.Where("room_id = ? AND state = ?", roomID, "none").
		Find(&noneCards).Error; err != nil {
		return fmt.Errorf("none 상태 카드 조회 중 에러 발생: %v", err)
	}
	// 1. room_id로 조회하여 state가 "opened"인 카드들의 state를 "none"으로 변경
	if err := tx.Model(&mysql.UserBirdCards{}).
		Where("room_id = ? AND state = ?", roomID, "opened").
		Update("state", "none").Error; err != nil {
		return fmt.Errorf("카드 상태 변경 (opened -> none) 중 에러 발생: %v", err)
	}

	// 3. none 상태 카드의 수를 확인
	noneCardCount := len(noneCards)
	// 4. none 상태 카드가 3개 미만이면 모두 "opened"로 변경
	if noneCardCount < 3 {
		var noneCardIDs []int
		for _, card := range noneCards {
			noneCardIDs = append(noneCardIDs, card.CardID)
		}

		if err := tx.Model(&mysql.UserBirdCards{}).
			Where("room_id = ? and id IN ?", itemChangeEntity.RoomID, noneCardIDs).
			Update("state", "opened").Error; err != nil {
			return fmt.Errorf("카드 상태 변경 (none -> opened) 중 에러 발생: %v", err)
		}
		return nil
	}

	// 5. none 상태 카드가 3개 이상이면 랜덤으로 3개 선택하여 "opened"로 변경
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
		Where("room_id = ? and card_id IN ?", itemChangeEntity.RoomID, selectedCardIDs).
		Update("state", "opened").Error; err != nil {
		return fmt.Errorf("카드 상태 변경 (none -> opened) 중 에러 발생: %v", err)
	}

	return nil
}
func ItemChangeConsumeUserItems(ctx context.Context, tx *gorm.DB, itemChangeEntity entity.WSItemChangeEntity) error {
	// 1. user_id, item_id로 조회하여 remaining_uses를 -1
	if err := tx.Model(&mysql.UserItems{}).
		Where("room_id = ? AND user_id = ? AND item_id = ?", itemChangeEntity.RoomID, itemChangeEntity.UserID, itemChangeEntity.ItemID).
		Update("remaining_uses", gorm.Expr("remaining_uses - 1")).Error; err != nil {
		return fmt.Errorf("user_item remaining_uses -1 중 에러 발생: %v", err)
	}

	return nil
}
