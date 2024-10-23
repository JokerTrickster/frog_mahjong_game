package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func RequestWinFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}
func RequestWinFindOneDora(c context.Context, tx *gorm.DB, roomID uint) (*mysql.Cards, error) {
	dora := mysql.Cards{}
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and state = ?", roomID, "dora").First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다. %v", err.Error())
	}
	return &dora, nil
}

// 카드 정보 체크 (소유하고 있는지 체크)
func RequestWinFindAllCards(c context.Context, tx *gorm.DB, requestWinEntity *entity.V2WSRequestWinEntity) ([]*mysql.Cards, error) {
	cards := make([]*mysql.Cards, 0)
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and user_id = ? and card_id IN ?", requestWinEntity.RoomID, requestWinEntity.UserID, requestWinEntity.Cards).Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다. %v", err.Error())
	}
	return cards, nil
}

// 카드 정보 모두 삭제
func RequestWinDeleteAllCards(ctx context.Context, tx *gorm.DB, requestWinEntity *entity.V2WSRequestWinEntity) error {
	err := tx.Model(&mysql.Cards{}).Where("room_id = ?", requestWinEntity.RoomID).Delete(&mysql.Cards{}).Error
	if err != nil {
		return fmt.Errorf("카드 삭제 실패 %v", err.Error())
	}
	return nil
}

// 유저 상태 변경 (play -> wait)
func RequestWinUpdateRoomUsers(c context.Context, tx *gorm.DB, requestWinEntity *entity.V2WSRequestWinEntity) error {
	err := tx.Model(&mysql.RoomUsers{}).Where("room_id = ?", requestWinEntity.RoomID).Update("player_state", "wait").Error
	if err != nil {
		return fmt.Errorf("방 유저 상태 변경 실패 %v", err.Error())
	}
	return nil
}
