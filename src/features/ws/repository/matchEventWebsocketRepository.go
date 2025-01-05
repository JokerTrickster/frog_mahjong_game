package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/features/ws/model/request"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MatchFindAllRoomUsers retrieves all room users with necessary preloads
func MatchFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Table("frog_room_users").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ?", roomID).
		Preload("User").
		Preload("Room").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).
		Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "room_users 조회 실패",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

// MatchFindOneRoomUsers retrieves the room ID for a specific user
func MatchFindOneRoomUsers(ctx context.Context, userID uint) (uint, *entity.ErrorInfo) {
	roomUser := mysql.FrogRoomUsers{}
	err := mysql.GormMysqlDB.WithContext(ctx).Where("user_id = ?", userID).First(&roomUser).Error
	if err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "방 유저 정보 조회 에러",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return uint(roomUser.RoomID), nil
}

// MatchFindOneWaitingRoom retrieves a waiting room that matches the criteria
func MatchFindOneWaitingRoom(ctx context.Context, count, timer uint) (*mysql.Rooms, *entity.ErrorInfo) {
	var room mysql.Rooms
	err := mysql.GormMysqlDB.Model(&mysql.Rooms{}).
		Where("min_count = ?", count).
		Where("max_count = ?", count).
		Where("timer = ?", timer).
		Where("state = ?", "wait").
		Where("current_count < max_count").
		Where("game_id = ?", 1).
		First(&room).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &mysql.Rooms{}, nil
		}
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "대기 방 조회 에러",
			Type: _errors.ErrFetchFailed,
		}
	}
	return &room, nil
}

// MatchFindOneAndUpdateUser updates a user's room and state
func MatchFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) *entity.ErrorInfo {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "ready",
	}
	err := tx.WithContext(ctx).Model(&mysql.Users{}).Where("id = ?", uID).Updates(user).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "유저 정보 업데이트 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// MatchInsertOneRoom inserts a new room
func MatchInsertOneRoom(ctx context.Context, room *mysql.Rooms) (int, *entity.ErrorInfo) {
	if ((room.MaxCount >= room.MinCount) && (room.MaxCount >= 2 || room.MinCount >= 2)) == false {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeBadRequest, // 400
			Msg:  "방 생성 조건이 유효하지 않습니다.",
			Type: _errors.ErrInvalidRequest,
		}
	}
	err := mysql.GormMysqlDB.WithContext(ctx).Create(&room).Error
	if err != nil {
		return 0, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "방 생성 실패",
			Type: _errors.ErrCreateFailed,
		}
	}
	return int(room.ID), nil
}

// MatchInsertOneRoomUser inserts a new room user
func MatchInsertOneRoomUser(ctx context.Context, tx *gorm.DB, roomUser *mysql.FrogRoomUsers) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Create(&roomUser).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "방 유저 생성 실패",
			Type: _errors.ErrCreateFailed,
		}
	}
	return nil
}

// MatchFindOneRoom retrieves a specific room by ID
func MatchFindOneRoom(ctx context.Context, tx *gorm.DB, req *request.ReqWSJoin) (mysql.Rooms, *entity.ErrorInfo) {
	room := mysql.Rooms{}
	err := tx.WithContext(ctx).Where("id = ?", req.RoomID).First(&room).Error
	if err != nil {
		return mysql.Rooms{}, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "방 정보를 찾을 수 없습니다.",
			Type: _errors.ErrRoomNotFound,
		}
	}
	return room, nil
}

// MatchFindOneAndUpdateRoom updates room's current player count
func MatchFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count + 1")).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "방 인원수 업데이트 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// MatchFindOneAndDeleteRoomUser deletes a specific room user
func MatchFindOneAndDeleteRoomUser(ctx context.Context, tx *gorm.DB, uID uint) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Where("user_id = ?", uID).Delete(&mysql.FrogRoomUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "방 유저 삭제 실패",
			Type: _errors.ErrDeleteRoomUserFailed,
		}
	}
	return nil
}

// MatchDeleteFrogCards deletes frog cards for a specific user
func MatchDeleteFrogCards(ctx context.Context, tx *gorm.DB, uID uint) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Model(&mysql.FrogUserCards{}).Where("user_id = ?", uID).Delete(&mysql.FrogUserCards{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "frog cards 삭제 실패",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}
func MatchPlayerStateUpdate(ctx context.Context, roomID, userID uint) *entity.ErrorInfo {
	err := mysql.GormMysqlDB.WithContext(ctx).Model(&mysql.FrogRoomUsers{}).
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
