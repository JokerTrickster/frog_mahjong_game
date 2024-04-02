package repository

import (
	"context"
	"errors"
	_interface "main/features/game/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewStartGameRepository(gormDB *gorm.DB) _interface.IStartGameRepository {
	return &StartGameRepository{GormDB: gormDB}
}

// 방장이 시작했는지 체크
func (g *StartGameRepository) CheckOwner(c context.Context, email string, roomID uint) error {
	room := mysql.Rooms{}
	err := g.GormDB.Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return err
	}
	if room.Owner != email {
		return errors.New("you are not the owner of this room")
	}
	return nil
}

// 방 유저들이 모두 준비했는지 체크
func (g *StartGameRepository) CheckReady(c context.Context, roomID uint) ([]mysql.RoomUsers, error) {

	roomUsers := make([]mysql.RoomUsers, 0)
	err := g.GormDB.Where("room_id = ?", roomID).Find(&roomUsers).Error
	if err != nil {
		return nil, err
	}

	return roomUsers, nil
}

// 방 유저들 상태 업데이트 (ready -> play)
func (g *StartGameRepository) UpdateRoomUser(c context.Context, roomID uint, state string) error {
	roomUsers := mysql.RoomUsers{
		PlayerState: state,
	}
	err := g.GormDB.Model(&roomUsers).Where("room_id = ? and player_state = ?", roomID, "ready").Updates(roomUsers).Error
	if err != nil {
		return err
	}
	return nil
}

// 방 상태 업데이트 (wait -> play)
func (g *StartGameRepository) UpdateRoom(c context.Context, roomID uint, state string) error {
	err := g.GormDB.Model(&mysql.Rooms{}).Where("id = ? and state = ?", roomID, "wait").Update("state", "play").Error
	if err != nil {
		return err
	}

	return nil
}

// 카드 데이터 생성
func (g *StartGameRepository) CreateCards(c context.Context, roomID uint, cards []mysql.Cards) error {
	err := g.GormDB.Create(&cards).Error
	if err != nil {
		return err
	}

	return nil
}
