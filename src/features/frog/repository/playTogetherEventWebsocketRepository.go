package repository

import (
	"context"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/features/ws/model/request"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PlayTogetherFindAllRoomUsers retrieves all room users with necessary preloads
func PlayTogetherFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
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

// PlayTogetherFindOneRoomUsers retrieves the room ID for a specific user
func PlayTogetherFindOneRoomUsers(ctx context.Context, userID uint) (uint, *entity.ErrorInfo) {
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

// PlayTogetherFindOneWaitingRoom retrieves a waiting room that matches the criteria
func PlayTogetherFindOneWaitingRoom(ctx context.Context, count, timer uint) (*mysql.Rooms, *entity.ErrorInfo) {
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
			Msg:  "대기 방 조회 중 에러 발생",
			Type: _errors.ErrFetchFailed,
		}
	}
	return &room, nil
}

// PlayTogetherFindOneAndUpdateUser updates a user's room and state
func PlayTogetherFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) *entity.ErrorInfo {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "play",
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

// PlayTogetherInsertOneRoom inserts a new room
func PlayTogetherInsertOneRoom(ctx context.Context, room mysql.Rooms) (int, *entity.ErrorInfo) {
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

// PlayTogetherInsertOneRoomUser inserts a new room user
func PlayTogetherInsertOneRoomUser(ctx context.Context, tx *gorm.DB, roomUser mysql.FrogRoomUsers) *entity.ErrorInfo {
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

// PlayTogetherFindOneRoom retrieves a specific room by ID
func PlayTogetherFindOneRoom(ctx context.Context, tx *gorm.DB, req *request.ReqWSJoin) (mysql.Rooms, *entity.ErrorInfo) {
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

// PlayTogetherFindOneAndUpdateRoom updates room details
func PlayTogetherFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID, count, timer uint) *entity.ErrorInfo {
	room := mysql.Rooms{
		MaxCount: int(count),
		MinCount: int(count),
		Timer:    int(timer),
	}
	err := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Updates(room).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "방 정보 업데이트 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// PlayTogetherAddPlayerToRoom increments the current player count in the room
func PlayTogetherAddPlayerToRoom(ctx context.Context, tx *gorm.DB, RoomID uint) *entity.ErrorInfo {
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

// PlayTogetherFindOneAndDeleteRoomUser deletes a specific room user
func PlayTogetherFindOneAndDeleteRoomUser(ctx context.Context, tx *gorm.DB, uID uint) *entity.ErrorInfo {
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
