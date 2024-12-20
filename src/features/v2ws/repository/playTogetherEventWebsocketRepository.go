package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/model/request"
	"main/utils"
	"main/utils/db/mysql"
	_redis "main/utils/db/redis"
	"time"

	"gorm.io/gorm"
)

func PlayTogetherFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("UserItems").Preload("Room").Preload("RoomMission").Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err)
	}
	return roomUsers, nil
}
func PlayTogetherFindOneRoomUsers(ctx context.Context, userID uint) (uint, error) {
	roomUser := mysql.RoomUsers{}
	err := mysql.GormMysqlDB.WithContext(ctx).Where("user_id = ?", userID).First(&roomUser).Error
	if err != nil {
		return 0, fmt.Errorf("방 유저 정보 조회 에러: %v", err)
	}
	return uint(roomUser.RoomID), nil
}

func PlayTogetherFindOneWaitingRoom(ctx context.Context, count, timer uint) (*mysql.Rooms, error) {
	var roomsDTO *mysql.Rooms
	err := mysql.GormMysqlDB.Model(&mysql.Rooms{}).Where("deleted_at is null and min_count = ? and max_count = ? and timer = ? and state = ? and current_count < max_count", count, count, timer, "wait").First(&roomsDTO).Error
	if err != nil {
		if err.Error() == "record not found" {
			return &mysql.Rooms{}, nil
		}
		return &mysql.Rooms{}, fmt.Errorf("대기 방 조회시 에러 발생: %v", err)
	}
	return roomsDTO, nil
}

func PlayTogetherFindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error {
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
func PlayTogetherInsertOneRoom(ctx context.Context, RoomDTO mysql.Rooms) (int, error) {
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
func PlayTogetherInsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomUserDTO mysql.RoomUsers) error {
	result := tx.WithContext(ctx).Create(&RoomUserDTO)
	if result.RowsAffected == 0 {
		return fmt.Errorf("방 유저 정보 생성 실패")
	}
	if result.Error != nil {
		return fmt.Errorf("방 유저 정보 생성 실패: %v", result.Error)
	}
	return nil
}

func PlayTogetherFindOneRoom(ctx context.Context, tx *gorm.DB, req *request.ReqWSJoin) (mysql.Rooms, error) {
	// 방 참여 가능한지 체크
	RoomDTO := mysql.Rooms{}
	result := tx.WithContext(ctx).Where("id = ?", req.RoomID).First(&RoomDTO)
	if result.Error != nil {
		return mysql.Rooms{}, fmt.Errorf("방 정보를 찾을 수 없습니다. %v", result.Error)
	}
	return RoomDTO, nil
}
func PlayTogetherFindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID, count, timer uint) error {
	room := mysql.Rooms{
		MaxCount: int(count),
		MinCount: int(count),
		Timer:    int(timer),
	}
	result := tx.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Updates(room)
	if result.Error != nil {
		return fmt.Errorf("방 정보 업데이트 실패: %v", result.Error)
	}

	return nil
}

func PlayTogetherAddPlayerToRoom(ctx context.Context, tx *gorm.DB, RoomID uint) error {
	result := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count + 1"))
	if result.Error != nil {
		return fmt.Errorf("방 인원수 업데이트 실패: %v", result.Error)
	}
	return nil
}

func PlayTogetherindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error {
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

func PlayTogetherFindOneAndDeleteRoomUser(ctx context.Context, tx *gorm.DB, uID uint) error {
	result := tx.WithContext(ctx).Where("user_id = ? ", uID).Delete(&mysql.RoomUsers{})
	// 방 유저 정보가 없는 경우
	if result.RowsAffected == 0 {
		return nil
	}
	if result.Error != nil {
		return fmt.Errorf("failed to delete room user: %v", result.Error)
	}
	return nil
}

func PlayTogetherDeleteRooms(ctx context.Context, uID uint) error {
	result := mysql.GormMysqlDB.WithContext(ctx).Where("owner_id = ?", uID).Delete(&mysql.Rooms{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete rooms: %v", result.Error)
	}
	return nil
}

func PlayTogetherDeleteRoomUsers(ctx context.Context, uID uint) error {
	result := mysql.GormMysqlDB.WithContext(ctx).Where("user_id = ?", uID).Delete(&mysql.RoomUsers{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete room users: %v", result.Error)
	}
	return nil
}

func PlayTogetherCreateMissions(ctx context.Context, tx *gorm.DB, roomID uint) error {
	// 랜덤으로 미션 ID 3개를 가져온다.
	var missionIDs []int
	err := tx.WithContext(ctx).
		Model(&mysql.Missions{}).
		Order("RAND()").
		Limit(3).
		Pluck("id", &missionIDs).Error
	if err != nil {
		return fmt.Errorf("미션 조회 실패: %v", err.Error())
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
		return fmt.Errorf("미션 생성 실패: %v", err.Error())
	}

	return nil
}

func PlayTogetherFindAllItems(ctx context.Context, tx *gorm.DB) ([]mysql.Items, error) {
	var items []mysql.Items
	err := tx.WithContext(ctx).Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("아이템 조회 실패: %v", err.Error())
	}
	return items, nil
}

func PlayTogetherInsertOneUserItem(ctx context.Context, tx *gorm.DB, userItemDTO mysql.UserItems) error {
	result := tx.WithContext(ctx).Create(&userItemDTO)
	if result.RowsAffected == 0 {
		return fmt.Errorf("유저 아이템 생성 실패")
	}
	if result.Error != nil {
		return fmt.Errorf("유저 아이템 생성 실패: %v", result.Error)
	}
	return nil
}
func PlayTogetherRedisSessionGet(ctx context.Context, sessionID string) (uint, error) {
	redisKey := fmt.Sprintf("abnormal_session:%s", sessionID)
	roomID, err := _redis.Client.Get(ctx, redisKey).Uint64()
	if err != nil {
		return 0, nil
	}
	return uint(roomID), nil
}

func PlayTogetherRedisSessionSet(ctx context.Context, sessionID string, roomID uint) error {
	redisKey := fmt.Sprintf("abnormal_session:%s", sessionID)
	err := _redis.Client.Set(ctx, redisKey, roomID, 3*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("세션 저장 실패: %v", err)
	}
	return nil
}
