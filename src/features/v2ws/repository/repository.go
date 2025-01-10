package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/utils/db/mysql"
	_redis "main/utils/db/redis"

	"gorm.io/gorm"
)

type JoinEventWebsocketRepository struct {
	GormDB *gorm.DB
}

func FindAllOpenCards(c context.Context, roomID int) ([]int, *entity.ErrorInfo) {
	var cards []int
	if err := mysql.GormMysqlDB.WithContext(c).Model(&mysql.UserBirdCards{}).Where("room_id = ? and state = ?", roomID, "opened").Pluck("card_id", &cards).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("FindAllOpenCards: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return cards, nil
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

func DeleteAllUserBirdCards(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.UserBirdCards{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllUserCards: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllRoomUsers(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.RoomUsers{}).Error
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
	err := tx.Where("owner_id = ?", userID).Delete(&mysql.Rooms{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllRooms: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllUserMissions(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.UserMissions{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllUserMissions: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllUserItems(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.UserItems{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllUserItems: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}
func DeleteAllUserMissionCards(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.UserMissionCards{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllUserMissionCards: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}
