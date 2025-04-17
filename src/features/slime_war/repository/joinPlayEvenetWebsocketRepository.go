package repository

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/utils/db/mysql"
	_redis "main/utils/db/redis"
	"time"

	"gorm.io/gorm"
)

func JoinPlayFindOneRoomUsers(ctx context.Context, userID uint) (uint, *entity.ErrorInfo) {
	roomUser := mysql.GameRoomUsers{}
	err := mysql.GormMysqlDB.WithContext(ctx).Where("user_id = ?", userID).First(&roomUser).Error
	if err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 정보 조회 에러: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return uint(roomUser.RoomID), nil
}

func JoinPlayFindOneWaitingRoom(ctx context.Context, password string) (*mysql.GameRooms, *entity.ErrorInfo) {
	var roomsDTO mysql.GameRooms
	err := mysql.GormMysqlDB.Model(&mysql.GameRooms{}).
		Where("deleted_at IS NULL").
		Where("password = ?", password).
		Where("state = ?", "wait").
		Where("current_count < max_count").
		Where("game_id = ?", 1).
		First(&roomsDTO).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &entity.ErrorInfo{
				Code: _errors.ErrCodeNotFound,
				Msg:  "비밀번호가 잘못되었습니다.",
				Type: _errors.ErrRoomNotFound,
			}
		}
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("대기 방 조회 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return &roomsDTO, nil
}

func JoinPlayFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) *entity.ErrorInfo {
	user := mysql.GameUsers{
		RoomID: int(RoomID),
		State:  "ready",
	}
	err := tx.WithContext(ctx).Model(&mysql.GameUsers{}).Where("id = ?", uID).Updates(user).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 정보 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func JoinPlayInsertOneRoom(ctx context.Context, RoomDTO mysql.Rooms) (int, *entity.ErrorInfo) {
	if !(RoomDTO.MaxCount >= RoomDTO.MinCount && (RoomDTO.MaxCount >= 2 || RoomDTO.MinCount >= 2)) {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeBadRequest,
			Msg:  "방 생성 조건이 올바르지 않습니다.",
			Type: _errors.ErrInvalidRequest,
		}
	}
	result := mysql.GormMysqlDB.WithContext(ctx).Create(&RoomDTO)
	if result.Error != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 생성 실패: %v", result.Error.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return int(RoomDTO.ID), nil
}

func JoinPlayInsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomUserDTO mysql.GameRoomUsers) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Create(&RoomUserDTO).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 생성 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func JoinPlayFindOneRoom(ctx context.Context, roomID uint) (*mysql.GameRooms, *entity.ErrorInfo) {
	var room mysql.GameRooms
	err := mysql.GormMysqlDB.WithContext(ctx).Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 정보 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomNotFound,
		}
	}
	return &room, nil
}

func JoinPlayFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Model(&mysql.GameRooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count + 1")).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 인원수 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func JoinPlayDeleteRoomUsers(ctx context.Context, userID uint) *entity.ErrorInfo {
	err := mysql.GormMysqlDB.WithContext(ctx).Where("user_id = ?", userID).Delete(&mysql.RoomUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 삭제 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func JoinRedisSessionGet(ctx context.Context, sessionID string) (uint, *entity.ErrorInfo) {
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

func JoinRedisSessionSet(ctx context.Context, sessionID string, roomID uint) *entity.ErrorInfo {
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

func JoinPlayDeleteRooms(ctx context.Context, userID uint) *entity.ErrorInfo {
	result := mysql.GormMysqlDB.WithContext(ctx).Where("owner_id = ?", userID).Delete(&mysql.Rooms{})
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 삭제 실패: %v", result.Error.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func JoinFindAllItems(ctx context.Context, tx *gorm.DB) ([]mysql.Items, *entity.ErrorInfo) {
	var items []mysql.Items
	err := tx.WithContext(ctx).Find(&items).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("아이템 조회 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return items, nil
}

func JoinInsertOneUserItem(ctx context.Context, tx *gorm.DB, userItemDTO mysql.UserItems) *entity.ErrorInfo {
	result := tx.WithContext(ctx).Create(&userItemDTO)
	if result.RowsAffected == 0 {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("아이템 생성 실패: %v", result.Error.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("아이템 생성 실패: %v", result.Error.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func JoinPlayPlayerStateUpdate(ctx context.Context, roomID, userID uint) *entity.ErrorInfo {
	err := mysql.GormMysqlDB.WithContext(ctx).Model(&mysql.RoomUsers{}).
		Where("room_id = ?", roomID).
		Where("user_id = ?", userID).
		Update("player_state", "play").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 상태 변경 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
