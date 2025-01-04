package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"time"

	_errors "main/features/ws/model/errors"
	_redis "main/utils/db/redis"

	"gorm.io/gorm"
)

type JoinEventWebsocketRepository struct {
	GormDB *gorm.DB
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
