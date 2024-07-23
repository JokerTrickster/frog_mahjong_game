package repository

import (
	"context"
	_errors "main/features/rooms/model/errors"
	_interface "main/features/rooms/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewOutRoomsRepository(gormDB *gorm.DB) _interface.IOutRoomsRepository {
	return &OutRoomsRepository{GormDB: gormDB}
}
func (g *OutRoomsRepository) FindOneUser(ctx context.Context, uID uint) (mysql.Users, error) {
	var user mysql.Users
	result := g.GormDB.WithContext(ctx).Where("id = ?", uID).First(&user)
	if result.Error != nil {
		return mysql.Users{}, utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)
	}
	return user, nil
}

func (g *OutRoomsRepository) ChangeRoomOnwer(ctx context.Context, RoomID uint, ownerID uint) error {
	var room mysql.Rooms
	result := g.GormDB.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Update("owner_id", ownerID)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)
	}
	return nil
}

func (g *OutRoomsRepository) FindOneRoomUser(ctx context.Context, RoomID uint) (mysql.RoomUsers, error) {
	var roomUser mysql.RoomUsers
	result := g.GormDB.WithContext(ctx).Where("room_id = ?", RoomID).First(&roomUser)
	if result.Error != nil {
		return mysql.RoomUsers{}, utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)
	}
	return roomUser, nil
}

// 방 삭제
func (g *OutRoomsRepository) FindOneAndDeleteRoom(ctx context.Context, RoomID uint) error {
	var room mysql.Rooms
	result := g.GormDB.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Delete(&room)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), _errors.ErrRoomNotFound.Error(), utils.ErrFromClient)
	}
	return nil
}

//

func (g *OutRoomsRepository) FindOneAndDeleteRoomUser(ctx context.Context, uID uint, RoomsID uint) error {
	var roomUser mysql.RoomUsers
	result := g.GormDB.WithContext(ctx).Model(&roomUser).Where("user_id = ? and room_id = ?", uID, RoomsID).Delete(&mysql.RoomUsers{})
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), _errors.ErrRoomUserNotFound.Error(), utils.ErrFromClient)
	}
	return nil
}

func (g *OutRoomsRepository) FindOneAndUpdateRoom(ctx context.Context, RoomID uint) (mysql.Rooms, error) {
	// 방 인원 -1
	var room mysql.Rooms
	result := g.GormDB.WithContext(ctx).Model(&room).Where("id = ?", RoomID).First(&room)
	if result.Error != nil {
		return mysql.Rooms{}, utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)
	}
	room.CurrentCount--
	result = g.GormDB.WithContext(ctx).Model(&room).Where("id = ?", RoomID).Updates(room)
	if result.Error != nil {
		return mysql.Rooms{}, utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)
	}

	return room, nil
}

func (g *OutRoomsRepository) FindOneAndUpdateUser(ctx context.Context, uID uint) error {
	user := mysql.Users{
		State:  "wait",
		RoomID: 1,
	}
	result := g.GormDB.WithContext(ctx).Model(&user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)

	}

	return nil
}
