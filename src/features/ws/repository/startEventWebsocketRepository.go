package repository

import (
	"context"
	"fmt"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

// 방장이 시작했는지 체크
func StartCheckOwner(ctx context.Context, tx *gorm.DB, uID uint, roomID uint) error {
	room := mysql.Rooms{}
	err := tx.WithContext(ctx).Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return fmt.Errorf("방 정보를 찾을 수 없습니다. %v", err)
	}
	if room.OwnerID != int(uID) {
		return fmt.Errorf("방장만 게임을 시작할 수 있습니다.")
	}
	return nil
}

// 방 유저들이 모두 준비했는지 체크
func StartCheckReady(ctx context.Context, tx *gorm.DB, roomID uint) ([]mysql.RoomUsers, error) {

	roomUsers := make([]mysql.RoomUsers, 0)
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).Find(&roomUsers).Error
	if err != nil {
		return nil, fmt.Errorf("방 유저 정보를 찾을 수 없습니다. %v", err)
	}

	return roomUsers, nil
}

// 방 유저 데이터 업데이트
func StartUpdateRoomUser(ctx context.Context, tx *gorm.DB, updateRoomUsers []mysql.RoomUsers) error {

	// 각 사용자 정보를 순회하며 각각 업데이트
	for _, user := range updateRoomUsers {
		err := tx.WithContext(ctx).Model(&mysql.RoomUsers{}).
			Where("room_id = ? AND user_id = ?", user.RoomID, user.UserID).
			Updates(user)

		if err.Error != nil {
			return fmt.Errorf("방 유저 정보 업데이트 실패: %v", err.Error)
		}
	}

	return nil
}

// 방 상태 업데이트 (wait -> play)
func StartUpdateRoom(ctx context.Context, tx *gorm.DB, roomID uint, state string) error {
	err := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ? and state = ?", roomID, "wait").Update("state", "play")
	if err.Error != nil {
		return fmt.Errorf("방 상태 업데이트 실패: %v", err.Error)
	}

	return nil
}

// 카드 데이터 생성
func StartCreateCards(ctx context.Context, tx *gorm.DB, roomID uint, cards []mysql.Cards) error {
	err := tx.WithContext(ctx).Create(&cards)
	if err.Error != nil {
		return fmt.Errorf("카드 정보 생성 실패: %v", err.Error)
	}

	return nil
}
