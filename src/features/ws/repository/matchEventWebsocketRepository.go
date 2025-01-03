package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/features/ws/model/request"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func MatchFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Table("frog_room_users").Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ?", roomID).
		Preload("User").
		Preload("Room").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 실패: %v", err.Error())
	}
	return roomUsers, nil
}

func MatchFindOneRoomUsers(ctx context.Context, userID uint) (uint, error) {
	roomUser := mysql.FrogRoomUsers{}
	err := mysql.GormMysqlDB.WithContext(ctx).Where("user_id = ?", userID).First(&roomUser).Error
	if err != nil {
		return 0, fmt.Errorf("방 유저 정보 조회 에러: %v", err)
	}
	return uint(roomUser.RoomID), nil
}

func MatchFindOneWaitingRoom(ctx context.Context, count, timer uint) (*mysql.Rooms, error) {
	var roomsDTO *mysql.Rooms
	err := mysql.GormMysqlDB.Model(&mysql.Rooms{}).
		Where("min_count = ?", count).
		Where("max_count = ?", count).
		Where("timer = ?", timer).
		Where("state = ?", "wait").
		Where("current_count < max_count").
		Where("game_id = ?", 1).
		First(&roomsDTO).Error
	if err != nil {
		if err.Error() == "record not found" {
			return &mysql.Rooms{}, nil
		}
		return &mysql.Rooms{}, fmt.Errorf("대기 방 조회시 에러 발생: %v", err)
	}
	return roomsDTO, nil
}

func MatchFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "ready",
	}
	result := tx.WithContext(ctx).Model(user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)
	}
	return nil
}
func MatchInsertOneRoom(ctx context.Context, RoomDTO *mysql.Rooms) (int, error) {
	//방 인원이 최대 인원이 최소 인원보다 많거나 같고, 최대 인원이 2명 이상이거나 최소 인원이 2명 이상이어야 한다.
	if ((RoomDTO.MaxCount >= RoomDTO.MinCount) && (RoomDTO.MaxCount >= 2 || RoomDTO.MinCount >= 2)) == false {
		return 0, utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), _errors.ErrBadRequest, utils.ErrFromClient)
	}
	result := mysql.GormMysqlDB.WithContext(ctx).Create(&RoomDTO)
	if result.RowsAffected == 0 {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), "failed room insert one", utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}
	return int(RoomDTO.ID), nil
}
func MatchInsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomUserDTO *mysql.FrogRoomUsers) error {
	result := tx.WithContext(ctx).Create(&RoomUserDTO)
	if result.RowsAffected == 0 {
		return fmt.Errorf("방 유저 정보 생성 실패")
	}
	if result.Error != nil {
		return fmt.Errorf("방 유저 정보 생성 실패: %v", result.Error)
	}
	return nil
}

func MatchFindOneRoom(ctx context.Context, tx *gorm.DB, req *request.ReqWSJoin) (mysql.Rooms, error) {
	// 방 참여 가능한지 체크
	RoomDTO := mysql.Rooms{}
	result := tx.WithContext(ctx).Where("id = ?", req.RoomID).First(&RoomDTO)
	if result.Error != nil {
		return mysql.Rooms{}, fmt.Errorf("방 정보를 찾을 수 없습니다. %v", result.Error)
	}
	return RoomDTO, nil
}

func MatchFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) error {
	result := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count + 1"))
	if result.Error != nil {
		return fmt.Errorf("방 인원수 업데이트 실패: %v", result.Error)
	}
	return nil
}

func MatchindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "play",
	}
	result := tx.WithContext(ctx).Model(&user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("유저 정보 업데이트 실패: %v", result.Error)
	}

	return nil
}

func MatchFindOneAndDeleteRoomUser(ctx context.Context, tx *gorm.DB, uID uint) error {
	result := tx.WithContext(ctx).Where("user_id = ? ", uID).Delete(&mysql.FrogRoomUsers{})
	// 방 유저 정보가 없는 경우
	if result.RowsAffected == 0 {
		return nil
	}
	if result.Error != nil {
		return fmt.Errorf("failed to delete room user: %v", result.Error)
	}
	return nil
}

func MatchDeleteFrogCards(ctx context.Context, tx *gorm.DB, uID uint) error {
	result := tx.WithContext(ctx).Model(&mysql.FrogUserCards{}).Where("user_id = ?", uID).Delete(&mysql.FrogUserCards{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete frog cards: %v", result.Error)
	}
	return nil
}
