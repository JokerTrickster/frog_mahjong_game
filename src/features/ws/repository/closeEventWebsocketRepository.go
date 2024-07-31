package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func CloseFindAllRoomUsers(ctx context.Context, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := mysql.GormMysqlDB.Preload("User").Preload("Room").Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err)
	}
	return roomUsers, nil
}
func CloseFindOneUser(ctx context.Context, tx *gorm.DB, uID uint) (mysql.Users, error) {
	var user mysql.Users
	result := tx.WithContext(ctx).Where("id = ?", uID).First(&user)
	if result.Error != nil {
		return mysql.Users{}, fmt.Errorf("유저 정보를 찾을 수 없습니다. %v", result.Error)
	}
	return user, nil
}

func CloseChangeRoomOnwer(ctx context.Context, tx *gorm.DB, RoomID uint, ownerID uint) error {
	var room mysql.Rooms
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Update("owner_id", ownerID)
	if result.Error != nil {
		return fmt.Errorf("방장 변경 실패: %v", result.Error)
	}
	return nil
}

func CloseFindOneRoomUser(ctx context.Context, tx *gorm.DB, RoomID uint) (mysql.RoomUsers, error) {
	var roomUser mysql.RoomUsers
	result := tx.WithContext(ctx).Where("room_id = ?", RoomID).First(&roomUser)
	if result.Error != nil {
		return mysql.RoomUsers{}, fmt.Errorf("방 유저 정보를 찾을 수 없습니다. %v", result.Error)
	}
	return roomUser, nil
}

// 방 삭제
func CloseFindOneAndDeleteRoom(ctx context.Context, tx *gorm.DB, RoomID uint) error {
	var room mysql.Rooms
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Delete(&room)
	if result.Error != nil {
		return fmt.Errorf("방 정보를 삭제할 수 없습니다. %v", result.Error)
	}
	return nil
}

//

func CloseFindOneAndDeleteRoomUser(ctx context.Context, tx *gorm.DB, uID uint, RoomsID uint) error {
	var roomUser mysql.RoomUsers
	result := tx.WithContext(ctx).Model(&roomUser).Where("user_id = ? and room_id = ?", uID, RoomsID).Delete(&mysql.RoomUsers{})
	if result.Error != nil {
		return fmt.Errorf("방 유저 정보를 삭제할 수 없습니다. %v", result.Error)
	}
	return nil
}

func CloseFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) (mysql.Rooms, error) {
	// 방 인원 -1
	var room mysql.Rooms
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", RoomID).First(&room)
	if result.Error != nil {
		return mysql.Rooms{}, fmt.Errorf("방 정보를 찾을 수 없습니다. %v", result.Error)
	}
	room.CurrentCount--
	result = tx.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Updates(room)
	if result.Error != nil {
		return mysql.Rooms{}, fmt.Errorf("방 인원을 업데이트할 수 없습니다. %v", result.Error)
	}

	return room, nil
}

func CloseFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint) error {
	user := mysql.Users{
		State:  "wait",
		RoomID: 1,
	}
	result := tx.WithContext(ctx).Model(&user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("유저 정보 업데이트 실패: %v", result.Error)

	}

	return nil
}
