package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func DoraCheckFirstPlayer(c context.Context, tx *gorm.DB, userID uint, roomID uint) error {
	var roomUsers mysql.RoomUsers
	err := tx.Model(&roomUsers).Where("user_id = ? AND room_id = ? and turn_number = ?", userID, roomID, 1).First(&roomUsers)
	if err.Error != nil {
		return fmt.Errorf("첫번째 플레이어가 아닙니다. %v", err.Error)
	}

	return nil
}

func DoraUpdateDoraCard(c context.Context, tx *gorm.DB, entity *entity.WSDoraEntity) error {

	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and card_id = ?", entity.RoomID, entity.CardID).Update("state", "dora")
	if err.Error != nil {
		return fmt.Errorf("도라 카드 업데이트 실패 %v", err.Error)
	}

	return nil
}
