package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func CancelMatchFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err)
	}
	return roomUsers, nil
}
func CancelMatchDeleteOneRoomUser(ctx context.Context, tx *gorm.DB, roomID, uID uint) error {
	roomUser := mysql.RoomUsers{}
	err := tx.Where("room_id = ? AND user_id = ?", roomID, uID).Delete(&roomUser).Error
	if err != nil {
		return err
	}
	return nil
}

func CancelMatchFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.Rooms, error) {
	var room mysql.Rooms
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).First(&room)
	if result.Error != nil {
		return &mysql.Rooms{}, fmt.Errorf("방 정보를 찾을 수 없습니다. %v", result.Error)
	}
	room.CurrentCount--
	result = tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).Updates(room)
	if result.Error != nil {
		return &mysql.Rooms{}, fmt.Errorf("방 인원을 업데이트할 수 없습니다. %v", result.Error)
	}

	if room.CurrentCount == 0 {
		result = tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).Delete(&room)
		if result.Error != nil {
			return &mysql.Rooms{}, fmt.Errorf("방 정보를 삭제할 수 없습니다. %v", result.Error)
		}
	}
	return &room, nil
}

func CancelMatchFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint) error {
	user := &mysql.Users{
		RoomID: 1,
		State:  "wait",
	}
	err := tx.Model(&user).Where("id = ?", uID).Updates(user).Error
	if err != nil {
		return err
	}

	return nil
}

func CancelMatchFindOneRoomUser(ctx context.Context, tx *gorm.DB, roomID uint) (uint, error) {
	var roomUser mysql.RoomUsers
	result := tx.WithContext(ctx).Where("room_id = ?", roomID).First(&roomUser)
	if result.Error != nil {
		return 0, fmt.Errorf("방 유저 정보를 찾을 수 없습니다. %v", result.Error)
	}
	return uint(roomUser.UserID), nil
}

func CancelMatchUpdateRoomOwner(ctx context.Context, tx *gorm.DB, roomID uint, roomUserID uint) error {
	var room mysql.Rooms
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", roomID).Update("owner_id", roomUserID)
	if result.Error != nil {
		return fmt.Errorf("방장 변경 실패: %v", result.Error)
	}
	return nil
}
