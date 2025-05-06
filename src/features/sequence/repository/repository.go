package repository

import (
	"context"
	"fmt"
	"main/features/sequence/model/entity"
	_errors "main/features/sequence/model/errors"
	"main/utils/db/mysql"
	_redis "main/utils/db/redis"

	"gorm.io/gorm"
)

type JoinEventWebsocketRepository struct {
	GormDB *gorm.DB
}

func PreloadUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.PreloadUsers, *entity.ErrorInfo) {
	var preloadUsers []entity.PreloadUsers

	if err := tx.Table("game_room_users").
		Preload("User").
		Preload("Room").
		Preload("SequenceRoomCards").
		Preload("SequenceRoomMaps").
		Preload("SequenceGameRoomSettings").
		Preload("SequenceUser").
		Where("room_id = ?", roomID).
		Find(&preloadUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room_users 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}

	return preloadUsers, nil
}

func ReconnectedUpdateRoomUser(c context.Context, roomID uint, userID uint) *entity.ErrorInfo {
	err := mysql.GormMysqlDB.Model(&mysql.RoomUsers{}).Where("room_id = ? and user_id = ?", roomID, userID).Update("player_state", "play").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("ReconnectedUpdateRoomUser: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func RedisSessionDelete(ctx context.Context, sessionID string) *entity.ErrorInfo {
	redisKey := fmt.Sprintf("abnormal_session:%s", sessionID)
	// Redis에서 키 삭제
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

func DeleteAllRoomUsers(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.GameRoomUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllRoomUsers: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllRooms(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("owner_id = ?", userID).Delete(&mysql.GameRooms{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllRooms: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllGameRooms(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("owner_id = ?", userID).Delete(&mysql.GameRooms{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllGameRooms: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllGameRoomUsers(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.GameRoomUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllGameRooms: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllSequenceUsers(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.SequenceUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllSequenceUsers: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}
