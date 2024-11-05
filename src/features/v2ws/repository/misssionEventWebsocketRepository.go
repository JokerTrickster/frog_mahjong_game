package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func MissionFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("RoomMission").Preload("BirdCards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Preload("UserMissions", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID)
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	fmt.Println("roomUsers: ", roomUsers[0].UserMissions)
	fmt.Println("roomUsers: ", roomUsers[1].UserMissions)
	return roomUsers, nil
}

// 카드 정보 체크 (소유하고 있는지 체크)
func MissionFindAllCards(c context.Context, tx *gorm.DB, missionEntity *entity.V2WSMissionEntity) error {
	cards := make([]*mysql.Cards, 0)
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and user_id = ? and card_id IN ?", missionEntity.RoomID, missionEntity.UserID, missionEntity.Cards).Find(&cards).Error
	if err != nil {
		return fmt.Errorf("카드를 찾을 수 없습니다. %v", err.Error())
	}
	return nil
}

func MissionCreateUserMission(ctx context.Context, userMissionDTO *mysql.UserMissions) (uint, error) {
	err := mysql.GormMysqlDB.Create(userMissionDTO).Error
	if err != nil {
		return 0, fmt.Errorf("유저 미션 생성 실패 %v", err.Error())
	}
	return userMissionDTO.ID, nil
}

func MissionCreateUserMissionCard(ctx context.Context, userMissionCardDTO *[]mysql.UserMissionCards) error {
	err := mysql.GormMysqlDB.Create(userMissionCardDTO).Error
	if err != nil {
		return fmt.Errorf("유저 미션 카드 생성 실패 %v", err.Error())
	}
	return nil
}
