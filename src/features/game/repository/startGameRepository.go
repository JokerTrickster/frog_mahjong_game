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

func (g *StartGameRepository) CheckReady(c context.Context, roomID uint) error {

	room := mysql.Room{}
	err := g.GormDB.Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return err
	}
	if room.State != "ready" {
		return errors.New("room is not ready")
	}
	return nil
}

func (g *StartGameRepository) UpdateRoomUser(c context.Context, roomID uint, state string) error {

	roomUsers := []mysql.RoomUser{}
	err := g.GormDB.Where("room_id = ?", roomID).Find(&roomUsers).Error
	if err != nil {
		return err
	}
	for _, ru := range roomUsers {
		ru.State = state
		err := g.GormDB.Save(&ru).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *StartGameRepository) UpdateRoom(c context.Context, roomID uint, state string) error {
	room := mysql.Room{}
	err := g.GormDB.Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return err
	}
	room.State = state
	err = g.GormDB.Save(&room).Error
	if err != nil {
		return err
	}
	return nil
}

func (g *StartGameRepository) CreateCards(c context.Context, roomID uint, cards []mysql.Cards) error {
	for _, card := range cards {
		err := g.GormDB.Create(&card).Error
		if err != nil {
			return err
		}
	}
	return nil
}
