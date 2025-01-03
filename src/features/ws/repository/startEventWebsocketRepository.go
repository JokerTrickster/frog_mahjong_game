package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StartFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Table("frog_room_users").Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ?", roomID).
		Preload("User").
		Preload("Room").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 실패: %v", err.Error())
	}
	return roomUsers, nil
}

func StartDeleteCards(ctx context.Context, tx *gorm.DB, roomID uint) error {
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).Delete(&mysql.FrogUserCards{})
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return fmt.Errorf("카드 삭제 실패 %v", err.Error)
	}
	return nil
}

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

func StartFindRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]mysql.FrogRoomUsers, error) {

	roomUsers := make([]mysql.FrogRoomUsers, 0)
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).Find(&roomUsers).Error
	if err != nil {
		return nil, fmt.Errorf("방 유저 정보를 찾을 수 없습니다. %v", err)
	}
	return roomUsers, nil
}

// 방 유저 데이터 업데이트
func StartUpdateRoomUser(ctx context.Context, tx *gorm.DB, updateRoomUsers []mysql.FrogRoomUsers) error {

	// 각 사용자 정보를 순회하며 각각 업데이트
	for _, user := range updateRoomUsers {
		err := tx.WithContext(ctx).Model(&mysql.FrogRoomUsers{}).
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
func StartCreateCards(ctx context.Context, tx *gorm.DB, cards []mysql.FrogUserCards) error {
	err := tx.WithContext(ctx).Create(&cards)
	if err.Error != nil {
		return fmt.Errorf("카드 정보 생성 실패: %v", err.Error)
	}

	return nil
}

func StartFindCards(ctx context.Context, tx *gorm.DB) ([]mysql.FrogCards, error) {
	var cards []mysql.FrogCards
	if err := tx.Find(&cards).Error; err != nil {
		return nil, fmt.Errorf("카드 조회 에러: %v", err.Error())
	}
	return cards, nil
}
