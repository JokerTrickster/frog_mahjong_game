package repository

import (
	"context"
	"fmt"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	_redis "main/utils/db/redis"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewCheckSessionRoomsRepository(gormDB *gorm.DB) _interface.ICheckSessionRoomsRepository {
	return &CheckSessionRoomsRepository{GormDB: gormDB}
}

func (d *CheckSessionRoomsRepository) RedisCheckSession(ctx context.Context, req *request.ReqCheckSession) (bool, error) {

	// redis에 세션 확인
	redisKey := fmt.Sprintf("abnormal_session:%s", req.SessionID)

	roomID, err := _redis.Client.Get(ctx, redisKey).Uint64()
	if err != nil {
		// Redis에서 key를 찾지 못한 경우 또는 다른 에러 발생
		if err == redis.Nil {
			return false, nil // Key가 없으므로 false 반환
		}
		return false, err // Redis 처리 중 다른 에러 발생
	}

	// roomID가 있으면 true 반환
	if roomID > 0 {
		return true, nil
	}

	// roomID가 0인 경우 false 반환
	return false, nil

	return true, nil
}
