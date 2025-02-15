package repository

import (
	"context"
	"fmt"
	"main/features/find_it/model/entity"
	_errors "main/features/find_it/model/errors"
	"main/utils/db/mysql"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func StartDeleteCards(ctx context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&mysql.UserBirdCards{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 삭제 실패: %v", err.Error()),
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

func StartCheckOwner(ctx context.Context, tx *gorm.DB, uID uint, roomID uint) *entity.ErrorInfo {
	room := mysql.GameRooms{}
	err := tx.WithContext(ctx).Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  fmt.Sprintf("방 정보를 찾을 수 없습니다. %v", err.Error()),
			Type: _errors.ErrRoomNotFound,
		}
	}
	if room.OwnerID != int(uID) {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeForbidden,
			Msg:  "방장만 게임을 시작할 수 있습니다.",
			Type: _errors.ErrUnauthorizedAction,
		}
	}
	return nil
}

func StartCheckReady(ctx context.Context, tx *gorm.DB, roomID uint) ([]mysql.RoomUsers, *entity.ErrorInfo) {
	roomUsers := make([]mysql.RoomUsers, 0)
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).Find(&roomUsers).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  fmt.Sprintf("방 유저 정보를 찾을 수 없습니다. %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

func StartUpdateRoomUser(ctx context.Context, tx *gorm.DB, updateRoomUsers []mysql.RoomUsers) *entity.ErrorInfo {
	for _, user := range updateRoomUsers {
		err := tx.WithContext(ctx).Model(&mysql.RoomUsers{}).
			Where("room_id = ? AND user_id = ?", user.RoomID, user.UserID).
			Updates(user).Error
		if err != nil {
			return &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("방 유저 정보 업데이트 실패: %v", err.Error()),
				Type: _errors.ErrUpdateFailed,
			}
		}
	}
	return nil
}

func StartUpdateRoom(ctx context.Context, tx *gorm.DB, roomID uint, roomUpdateData mysql.GameRooms) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Model(&mysql.GameRooms{}).
		Where("id = ? AND state = ?", roomID, "wait").
		Updates(roomUpdateData).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 상태 업데이트 실패: %v", err.Error()),
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

func StartCreateCards(ctx context.Context, tx *gorm.DB, cards []mysql.UserBirdCards) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Create(&cards).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 정보 생성 실패: %v", err.Error()),
			Type: _errors.ErrCreateFailed,
		}
	}
	return nil
}

func StartUpdateCardState(ctx context.Context, roomID uint) ([]int, *entity.ErrorInfo) {
	var birdCards []*mysql.BirdCards
	err := mysql.GormMysqlDB.Model(&mysql.BirdCards{}).Find(&birdCards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 총 수 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}

	openCards := make([]int, len(birdCards))
	for i := 0; i < len(birdCards); i++ {
		openCards[i] = int(birdCards[i].ID)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(birdCards), func(i, j int) {
		openCards[i], openCards[j] = openCards[j], openCards[i]
	})

	for i := 0; i < 3; i++ {
		err = mysql.GormMysqlDB.Model(&mysql.UserBirdCards{}).
			Where("card_id = ?", openCards[i]).
			Update("state", "opened").Error
		if err != nil {
			return nil, &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("카드 상태 업데이트 실패: %v", err.Error()),
				Type: _errors.ErrUpdateFailed,
			}
		}
	}

	return openCards[:3], nil
}

func StartCreateMissions(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
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
			Type: _errors.ErrFetchFailed,
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
			Type: _errors.ErrCreateFailed,
		}
	}
	return nil
}

func StartBirdCard(ctx context.Context, tx *gorm.DB) ([]*mysql.BirdCards, *entity.ErrorInfo) {
	var birdCards []*mysql.BirdCards
	err := tx.WithContext(ctx).Find(&birdCards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("birdCards 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return birdCards, nil
}

func StartCheckRoomState(ctx context.Context, roomID uint) (string, *entity.ErrorInfo) {
	room := mysql.GameRooms{}
	err := mysql.GormMysqlDB.WithContext(ctx).Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return "", &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound,
			Msg:  fmt.Sprintf("방 정보를 찾을 수 없습니다. %v", err.Error()),
			Type: _errors.ErrRoomNotFound,
		}
	}
	return room.State, nil
}
