package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	"main/utils/db/mysql"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func StartFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("RoomMission").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Preload("UserMissions", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID)
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}

func StartDeleteCards(ctx context.Context, tx *gorm.DB, userID uint) error {
	err := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&mysql.UserBirdCards{})
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return fmt.Errorf("카드 삭제 실패 %v", err.Error)
	}
	return nil
}

// 방장이 시작했는지 체크
func StartCheckOwner(ctx context.Context, tx *gorm.DB, uID uint, roomID uint) (uint, error) {
	room := mysql.Rooms{}
	err := tx.WithContext(ctx).Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return 0, fmt.Errorf("방 정보를 찾을 수 없습니다. %v", err)
	}
	if room.OwnerID != int(uID) {
		return 0, fmt.Errorf("방장만 게임을 시작할 수 있습니다.")
	}
	return uint(room.OwnerID), nil
}

// 방 유저들이 모두 준비했는지 체크
func StartCheckReady(ctx context.Context, tx *gorm.DB, roomID uint) ([]mysql.RoomUsers, error) {

	roomUsers := make([]mysql.RoomUsers, 0)
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).Find(&roomUsers).Error
	if err != nil {
		return nil, fmt.Errorf("방 유저 정보를 찾을 수 없습니다. %v", err)
	}

	return roomUsers, nil
}

// 방 유저 데이터 업데이트
func StartUpdateRoomUser(ctx context.Context, tx *gorm.DB, updateRoomUsers []mysql.RoomUsers) error {

	// 각 사용자 정보를 순회하며 각각 업데이트
	for _, user := range updateRoomUsers {
		err := tx.WithContext(ctx).Model(&mysql.RoomUsers{}).
			Where("room_id = ? AND user_id = ?", user.RoomID, user.UserID).
			Updates(user)

		if err.Error != nil {
			return fmt.Errorf("방 유저 정보 업데이트 실패: %v", err.Error)
		}
	}

	return nil
}

// 방 상태 업데이트 (wait -> play)
func StartUpdateRoom(ctx context.Context, tx *gorm.DB, roomID uint, state string) error {
	err := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ? and state = ?", roomID, "wait").Update("state", "play")
	if err.Error != nil {
		return fmt.Errorf("방 상태 업데이트 실패: %v", err.Error)
	}

	return nil
}

// 유저들 코인 -1 차감한다.
func StartDiffCoin(ctx context.Context, tx *gorm.DB, roomID uint) error {
	err := tx.WithContext(ctx).Model(&mysql.Users{}).Where("room_id = ?", roomID).Update("coin", gorm.Expr("coin - 1"))
	if err.Error != nil {
		return fmt.Errorf("유저 코인 차감 실패: %v", err.Error)
	}

	return nil
}

// 카드 데이터 생성
func StartCreateCards(ctx context.Context, tx *gorm.DB, cards []mysql.UserBirdCards) error {
	err := tx.WithContext(ctx).Create(&cards)
	if err.Error != nil {
		return fmt.Errorf("카드 정보 생성 실패: %v", err.Error)
	}

	return nil
}

// 랜덤으로 카드 3장 상태를 opened으로 변경한다.
func StartUpdateCardState(ctx context.Context, roomID uint) ([]int, error) {
	// 카드 총 수를 가져온다.
	var count int64
	err := mysql.GormMysqlDB.Model(&mysql.BirdCards{}).Count(&count).Error
	if err != nil {
		return nil, fmt.Errorf("카드 총 수 조회 실패: %v", err.Error())
	}

	// 카드 총 수에서 랜덤으로 3개의 카드 ID를 가져온다. (rand를 통해 3개의 카드를 뽑는다.)
	openCards := make([]int, count)
	for i := 0; i < int(count); i++ {
		openCards[i] = i + 1
	}

	// 배열을 랜덤하게 섞음
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(int(count), func(i, j int) {
		openCards[i], openCards[j] = openCards[j], openCards[i]
	})
	// opened 상태로 카드 상태 업데이트 한다.
	for i := 0; i < 3; i++ {
		err = mysql.GormMysqlDB.Model(&mysql.UserBirdCards{}).Where("card_id = ?", openCards[i]).Update("state", "opened").Error
		if err != nil {
			return nil, fmt.Errorf("카드 상태 업데이트 실패: %v", err.Error())
		}
	}

	return openCards[:3], nil
}

// 미션을 랜덤으로 3개 생성한다.
func StartCreateMissions(ctx context.Context, tx *gorm.DB, roomID uint) error {
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

func StartBirdCard(ctx context.Context, tx *gorm.DB) ([]*mysql.BirdCards, error) {
	var birdCards []*mysql.BirdCards
	err := tx.WithContext(ctx).Find(&birdCards).Error
	if err != nil {
		return nil, fmt.Errorf("birdCards 조회 실패: %v", err.Error())
	}

	return birdCards, nil
}
