package repository

import (
	"context"
	"errors"
	"fmt"
	"main/features/ws/model/entity"
	"time"

	_errors "main/features/ws/model/errors"
	"main/utils/db/mysql"
	_redis "main/utils/db/redis"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JoinEventWebsocketRepository struct {
	GormDB *gorm.DB
}

func PreloadFindGameInfo(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Table("frog_room_users").Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ?", roomID).
		Preload("User").
		Preload("Room").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil,
			&entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("room_users 조회 실패: %v", err.Error()),
				Type: _errors.ErrInternalServer,
			}
	}
	return roomUsers, nil
}
func FindOneDoraCard(ctx context.Context, roomID int) (*mysql.FrogUserCards, *entity.ErrorInfo) {
	doraCard := &mysql.FrogUserCards{}
	result := mysql.GormMysqlDB.Table("frog_user_cards").
		Where("room_id = ?", roomID).
		Where("state = ?", "dora").
		First(&doraCard)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil,
			&entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("도라카드 조회 실패: %v", result.Error),
				Type: _errors.ErrInternalServer,
			}
	}
	return doraCard, nil
}

func RedisSessionSet(ctx context.Context, sessionID string, roomID uint) *entity.ErrorInfo {
	redisKey := fmt.Sprintf("abnormal_session:%s", sessionID)
	err := _redis.Client.Set(ctx, redisKey, roomID, 3*time.Minute).Err()
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("세션 저장 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func MatchRedisSessionGet(ctx context.Context, sessionID string) (uint, *entity.ErrorInfo) {
	redisKey := fmt.Sprintf("abnormal_session:%s", sessionID)
	roomID, err := _redis.Client.Get(ctx, redisKey).Uint64()
	if err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("세션 조회 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return uint(roomID), nil
}

func RedisSessionDelete(ctx context.Context, sessionID string) *entity.ErrorInfo {
	redisKey := fmt.Sprintf("abnormal_session:%s", sessionID)
	err := _redis.Client.Del(ctx, redisKey).Err()
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("세션 삭제 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllRooms(ctx context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	if err := tx.Model(&mysql.Rooms{}).Where("owner_id = ?", userID).Delete(&mysql.Rooms{}).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 삭제 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllFrogUserCards(ctx context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	if err := tx.Model(&mysql.FrogUserCards{}).Where("user_id = ?", userID).Delete(&mysql.FrogUserCards{}).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 카드 삭제 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllFrogRoomUsers(ctx context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	if err := tx.Model(&mysql.FrogRoomUsers{}).Where("user_id = ?", userID).Delete(&mysql.FrogRoomUsers{}).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("룸 유저 삭제 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}
