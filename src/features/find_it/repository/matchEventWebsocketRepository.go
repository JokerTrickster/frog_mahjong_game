package repository

import (
	"context"
	"errors"
	"fmt"
	"main/features/find_it/model/entity"
	_errors "main/features/find_it/model/errors"
	"main/utils/db/mysql"
	_redis "main/utils/db/redis"
	"time"

	"gorm.io/gorm"
)

func MatchFindOneRoomUsers(ctx context.Context, userID uint) (uint, *entity.ErrorInfo) {
	roomUser := mysql.GameRoomUsers{}
	err := mysql.GormMysqlDB.WithContext(ctx).Where("user_id = ?", userID).First(&roomUser).Error
	if err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 정보 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return uint(roomUser.RoomID), nil
}
func MatchFindOneWaitingRoom(ctx context.Context) (*mysql.GameRooms, *entity.ErrorInfo) {
	var roomsDTO mysql.GameRooms

	query := mysql.GormMysqlDB.Model(&mysql.GameRooms{}).
		Where("deleted_at IS NULL").
		Where("state = ?", "wait").
		Where("current_count < max_count").
		Where("game_id = ?", 1)

	err := query.First(&roomsDTO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &entity.ErrorInfo{
				Code: _errors.ErrCodeNotFound,
				Msg:  err.Error(),
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
func MatchInsertOneRoomSetting(ctx context.Context, tx *gorm.DB, roomSettingDTO *mysql.FindItRoomSettings) *entity.ErrorInfo {
	result := tx.WithContext(ctx).Create(&roomSettingDTO)
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 설정 생성 실패: %v", result.Error.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func MatchFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) *entity.ErrorInfo {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "ready",
	}
	err := tx.WithContext(ctx).Model(&mysql.Users{}).Where("id = ?", uID).Updates(user).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 정보 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func MatchInsertOneRoom(ctx context.Context, RoomDTO *mysql.GameRooms) (int, *entity.ErrorInfo) {
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

func MatchInsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomUserDTO *mysql.GameRoomUsers) *entity.ErrorInfo {
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

func MatchFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) *entity.ErrorInfo {
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

func MatchCreateMissions(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	var missionIDs []int
	err := tx.WithContext(ctx).
		Model(&mysql.Missions{}).
		Order("RAND()").
		Limit(3).
		Pluck("id", &missionIDs).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("미션 조회 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}

	roomMissions := make([]mysql.RoomMissions, 0)
	for _, missionID := range missionIDs {
		roomMissions = append(roomMissions, mysql.RoomMissions{
			RoomID:    int(roomID),
			MissionID: missionID,
		})
	}
	err = tx.WithContext(ctx).Create(&roomMissions).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("미션 생성 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}

	return nil
}

func MatchFindOneRoomMission(ctx context.Context, tx *gorm.DB, roomID uint) ([]mysql.RoomMissions, *entity.ErrorInfo) {
	var roomMissions []mysql.RoomMissions
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).Find(&roomMissions).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 미션 조회 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return roomMissions, nil
}

func MatchFindAllItems(ctx context.Context, tx *gorm.DB) ([]mysql.Items, *entity.ErrorInfo) {
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

func MatchInsertOneUserItem(ctx context.Context, tx *gorm.DB, userItemDTO mysql.UserItems) *entity.ErrorInfo {
	result := tx.WithContext(ctx).Create(&userItemDTO)
	if result.RowsAffected == 0 {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "유저 아이템 생성 실패",
			Type: _errors.ErrInternalServer,
		}
	}
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 아이템 생성 실패: %v", result.Error.Error()),
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

func MatchRedisSessionSet(ctx context.Context, sessionID string, roomID uint) *entity.ErrorInfo {
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

func MatchDeleteRooms(ctx context.Context, uID uint) *entity.ErrorInfo {
	result := mysql.GormMysqlDB.WithContext(ctx).Where("owner_id = ?", uID).Delete(&mysql.Rooms{})
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 삭제 실패: %v", result.Error.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func MatchDeleteRoomUsers(ctx context.Context, uID uint) *entity.ErrorInfo {
	result := mysql.GormMysqlDB.WithContext(ctx).Where("user_id = ?", uID).Delete(&mysql.RoomUsers{})
	if result.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 삭제 실패: %v", result.Error.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func MatchPlayerStateUpdate(ctx context.Context, roomID, userID uint) *entity.ErrorInfo {
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
