package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/model/request"
	"main/utils/db/mysql"
	_redis "main/utils/db/redis"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func PlayTogetherFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("User").
		Preload("UserItems").
		Preload("Room").
		Preload("RoomMission").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).
		Preload("UserMissions", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID)
		}).
		Where("room_id = ?", roomID).
		Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room_users 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

func PlayTogetherFindOneRoomUsers(ctx context.Context, userID uint) (uint, *entity.ErrorInfo) {
	roomUser := mysql.RoomUsers{}
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

func PlayTogetherFindOneWaitingRoom(ctx context.Context, count, timer uint) (*mysql.Rooms, *entity.ErrorInfo) {
	var roomsDTO mysql.Rooms
	err := mysql.GormMysqlDB.Model(&mysql.Rooms{}).
		Where("deleted_at IS NULL AND min_count = ? AND max_count = ? AND timer = ? AND state = ? AND current_count < max_count", count, count, timer, "wait").
		First(&roomsDTO).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &entity.ErrorInfo{
				Code: _errors.ErrCodeNotFound,
				Msg:  "대기 중인 방을 찾을 수 없습니다.",
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

func PlayTogetherFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) *entity.ErrorInfo {
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

func PlayTogetherInsertOneRoom(ctx context.Context, RoomDTO mysql.Rooms) (int, *entity.ErrorInfo) {
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

func PlayTogetherInsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomUserDTO mysql.RoomUsers) *entity.ErrorInfo {
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

func PlayTogetherFindOneRoom(ctx context.Context, tx *gorm.DB, req *request.ReqWSJoin) (*mysql.Rooms, *entity.ErrorInfo) {
	var room mysql.Rooms
	err := tx.WithContext(ctx).Where("id = ?", req.RoomID).First(&room).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 정보 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomNotFound,
		}
	}
	return &room, nil
}

func PlayTogetherFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID, count, timer uint) *entity.ErrorInfo {
	room := mysql.Rooms{
		MaxCount: int(count),
		MinCount: int(count),
		Timer:    int(timer),
	}
	err := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Updates(room).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 정보 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func PlayTogetherAddPlayerToRoom(ctx context.Context, tx *gorm.DB, RoomID uint) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count + 1")).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 인원수 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func PlayTogetherDeleteRoomUsers(ctx context.Context, uID uint) *entity.ErrorInfo {
	err := mysql.GormMysqlDB.WithContext(ctx).Where("user_id = ?", uID).Delete(&mysql.RoomUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 삭제 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func PlayTogetherRedisSessionSet(ctx context.Context, sessionID string, roomID uint) *entity.ErrorInfo {
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

func PlayTogetherFindAllItems(ctx context.Context, tx *gorm.DB) ([]mysql.Items, *entity.ErrorInfo) {
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

func PlayTogetherInsertOneUserItem(ctx context.Context, tx *gorm.DB, userItemDTO mysql.UserItems) *entity.ErrorInfo {
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

func PlayTogetherDeleteRooms(ctx context.Context, uID uint) *entity.ErrorInfo {
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
func PlayTogetherCreateMissions(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	// 랜덤으로 미션 ID 3개를 가져온다.
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

	// 미션 정보를 생성한다.
	roomMissions := make([]mysql.RoomMissions, 0)
	for _, missionID := range missionIDs {
		roomMission := mysql.RoomMissions{
			RoomID:    int(roomID),
			MissionID: missionID,
		}
		roomMissions = append(roomMissions, roomMission)
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
