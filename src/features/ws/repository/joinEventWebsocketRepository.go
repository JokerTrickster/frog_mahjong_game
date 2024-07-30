package repository

import (
	"context"
	"errors"
	"log"
	"main/features/ws/model/entity"
	"main/features/ws/model/request"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func JoinFindAllRoomUsers(ctx context.Context, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := mysql.GormMysqlDB.Preload("User").Preload("Room").Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		log.Fatalf("RoomUsers 조회 에러: %s", err)
	}
	return roomUsers, nil
}

func JoinFindOneRoom(ctx context.Context, tx *gorm.DB, req *request.ReqWSJoin) (mysql.Rooms, error) {
	// 방 참여 가능한지 체크
	RoomDTO := mysql.Rooms{}
	result := tx.WithContext(ctx).Where("id = ? and password = ?", req.RoomID, req.Password).First(&RoomDTO)
	if result.Error != nil {
		return mysql.Rooms{}, errors.New("방이 존재하지 않습니다.")
	}
	return RoomDTO, nil
}

func JoinFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) error {
	result := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count + 1"))
	if result.Error != nil {
		return errors.New("방 인원 증가 실패")
	}
	return nil
}

func JoinFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "play",
	}
	result := tx.WithContext(ctx).Model(&user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return errors.New("유저 정보 업데이트 실패")
	}

	return nil
}

func JoinInsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomUserDTO mysql.RoomUsers) error {
	result := tx.WithContext(ctx).Create(&RoomUserDTO)
	if result.RowsAffected == 0 {
		return errors.New("방 유저 정보 생성 실패")
	}
	if result.Error != nil {
		return errors.New("방 유저 정보 생성 실패")
	}
	return nil
}
