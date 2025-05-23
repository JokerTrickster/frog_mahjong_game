package repository

import (
	"context"
	_errors "main/features/game/model/errors"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewStartGameRepository(gormDB *gorm.DB) _interface.IStartGameRepository {
	return &StartGameRepository{GormDB: gormDB}
}

// 방장이 시작했는지 체크
func (g *StartGameRepository) CheckOwner(c context.Context, uID uint, roomID uint) error {
	room := mysql.Rooms{}
	err := g.GormDB.Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrBadParameter, utils.Trace(), _errors.ErrBadRequest.Error(), utils.ErrFromClient)
	}
	if room.OwnerID != int(uID) {
		return utils.ErrorMsg(c, utils.ErrBadParameter, utils.Trace(), _errors.ErrNotOwner.Error(), utils.ErrFromClient)
	}
	return nil
}

// 방 유저들이 모두 준비했는지 체크
func (g *StartGameRepository) CheckReady(c context.Context, roomID uint) ([]mysql.RoomUsers, error) {

	roomUsers := make([]mysql.RoomUsers, 0)
	err := g.GormDB.Where("room_id = ?", roomID).Find(&roomUsers).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrBadParameter, utils.Trace(), _errors.ErrNotAllUsersReady.Error(), utils.ErrFromClient)
	}

	return roomUsers, nil
}

// 방 유저 데이터 업데이트
func (g *StartGameRepository) UpdateRoomUser(c context.Context, updateRoomUsers []mysql.RoomUsers) error {

	// 각 사용자 정보를 순회하며 각각 업데이트
	for _, user := range updateRoomUsers {
		err := g.GormDB.Model(&mysql.RoomUsers{}).
			Where("room_id = ? AND user_id = ?", user.RoomID, user.UserID).
			Updates(user)

		if err.Error != nil {
			return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error.Error(), utils.ErrFromMysqlDB)
		}
	}

	return nil
}

// 방 상태 업데이트 (wait -> play)
func (g *StartGameRepository) UpdateRoom(c context.Context, roomID uint, state string) error {
	err := g.GormDB.Model(&mysql.Rooms{}).Where("id = ? and state = ?", roomID, "wait").Update("state", "play")
	if err.Error != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error.Error(), utils.ErrFromMysqlDB)
	}

	return nil
}

// 카드 데이터 생성
func (g *StartGameRepository) CreateCards(c context.Context, roomID uint, cards []mysql.Cards) error {
	err := g.GormDB.Create(&cards)
	if err.Error != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error.Error(), utils.ErrFromMysqlDB)
	}

	return nil
}
