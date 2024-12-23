package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func MissionFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	err := tx.Preload("User").
		Preload("UserItems").
		Preload("Room").
		Preload("RoomMission").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Order("updated_at ASC")
		}).
		Preload("UserMissions", "room_id = ?", roomID).
		Where("room_id = ?", roomID).
		Find(&roomUsers).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room_users 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

// 카드 정보 체크 (소유하고 있는지 체크)
func MissionFindAllCards(ctx context.Context, tx *gorm.DB, missionEntity *entity.V2WSMissionEntity) *entity.ErrorInfo {
	cards := make([]*mysql.Cards, 0)
	err := tx.Model(&mysql.Cards{}).
		Where("room_id = ? AND user_id = ? AND card_id IN ?", missionEntity.RoomID, missionEntity.UserID, missionEntity.Cards).
		Find(&cards).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 조회 실패: %v", err.Error()),
			Type: _errors.ErrNotFoundCard,
		}
	}
	return nil
}

func MissionCreateUserMission(ctx context.Context, tx *gorm.DB, userMissionDTO *mysql.UserMissions) (uint, *entity.ErrorInfo) {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Create(userMissionDTO).Error
	if err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 미션 생성 실패: %v", err.Error()),
			Type: _errors.ErrCreateFailed,
		}
	}
	return userMissionDTO.ID, nil
}

func MissionCreateUserMissionCard(ctx context.Context, tx *gorm.DB, userMissionCardDTO *[]mysql.UserMissionCards) *entity.ErrorInfo {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Create(userMissionCardDTO).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 미션 카드 생성 실패: %v", err.Error()),
			Type: _errors.ErrCreateFailed,
		}
	}
	return nil
}
