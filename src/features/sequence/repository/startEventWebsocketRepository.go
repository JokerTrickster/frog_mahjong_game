package repository

import (
	"context"
	"fmt"
	"main/features/sequence/model/entity"
	_errors "main/features/sequence/model/errors"
	"main/utils/db/mysql"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func CreateRoundImages(ctx context.Context, tx *gorm.DB, roundImages []*mysql.FindItRoundImages) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Create(&roundImages).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("라운드 이미지 생성 실패: %v", err.Error()),
			Type: _errors.ErrCreateFailed,
		}
	}
	return nil
}

func FindImages(ctx context.Context, tx *gorm.DB) ([]*mysql.FindItImages, *entity.ErrorInfo) {
	var images []*mysql.FindItImages
	err := tx.WithContext(ctx).
		Order("RAND()").
		Limit(20).
		Find(&images).Error

	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("images 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return images, nil
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

func StartCreateSequenceCards(ctx context.Context, tx *gorm.DB, cards []mysql.SequenceRoomCards) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Create(&cards).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 생성 실패: %v", err.Error()),
			Type: _errors.ErrCreateFailed,
		}
	}
	return nil
}

func StartCreateSequenceMaps(ctx context.Context, tx *gorm.DB, maps []mysql.SequenceRoomMaps) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Create(&maps).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("맵 생성 실패: %v", err.Error()),
			Type: _errors.ErrCreateFailed,
		}
	}
	return nil
}

func StartUpdateSequenceUser(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	var users []mysql.SequenceUsers

	// roomID로 해당 유저 두 명 조회
	if err := tx.Where("room_id = ?", roomID).Find(&users).Error; err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("슬라임워 유저 조회 실패: %v", err.Error()),
			Type: _errors.ErrSequenceUsersNotFound,
		}
	}

	if len(users) != 2 {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInvalidUserCount,
			Msg:  "슬라임워 유저 수는 반드시 2명이어야 합니다.",
			Type: _errors.ErrSequenceUsersNotFound,
		}
	}

	// 랜덤으로 0 또는 1 배정
	rand.Seed(time.Now().UnixNano())
	r0 := rand.Intn(2) // 0 또는 1
	r1 := 1 - r0       // 반대값

	// 무작위 순서로 섞기
	if rand.Intn(2) == 0 {
		users[0].Turn = r0
		users[0].ColorType = r0
		users[1].Turn = r1
		users[1].ColorType = r1
	} else {
		users[0].Turn = r1
		users[0].ColorType = r1
		users[1].Turn = r0
		users[1].ColorType = r0
	}

	// DB 업데이트
	for _, u := range users {
		if err := tx.Model(&mysql.SequenceUsers{}).
			Where("id = ?", u.ID).
			Updates(map[string]interface{}{
				"turn":       u.Turn,
				"color_type": u.ColorType,
			}).Error; err != nil {
			return &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("슬라임워 유저 업데이트 실패: %v", err.Error()),
				Type: _errors.ErrSequenceUserUpdateFailed,
			}
		}
	}

	return nil
}

func StartFindRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]mysql.GameRoomUsers, *entity.ErrorInfo) {
	var roomUsers []mysql.GameRoomUsers
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).Find(&roomUsers).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 조회 실패: %v", err.Error()),
			Type: _errors.ErrFetchFailed,
		}
	}
	return roomUsers, nil
}

func StartCreateSequenceUserCards(ctx context.Context, tx *gorm.DB, roomUsers []mysql.GameRoomUsers) *entity.ErrorInfo {
	// Create a map to track used numbers
	used := make(map[int]bool)
	numbers := make([]int, 0, 10)

	// Generate 10 unique random numbers between 1 and 48
	for len(numbers) < 14 {
		num := rand.Intn(96) + 1
		if !used[num] {
			used[num] = true
			numbers = append(numbers, num)
		}
	}

	// 각 유저에게 5개의 카드 할당
	for i, user := range roomUsers {
		startIdx := i * 7
		userCards := numbers[startIdx : startIdx+5]

		//userCards에 해당되는 card_id에 user_id와 state를 업데이트한다.
		if err := tx.Model(&mysql.SequenceRoomCards{}).
			Where("card_id IN (?)", userCards).
			Updates(map[string]interface{}{
				"user_id": user.UserID,
				"state":   "owned",
			}).Error; err != nil {
			return &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal,
				Msg:  fmt.Sprintf("유저 카드 업데이트 실패: %v", err.Error()),
				Type: _errors.ErrUpdateFailed,
			}
		}
	}

	return nil
}
