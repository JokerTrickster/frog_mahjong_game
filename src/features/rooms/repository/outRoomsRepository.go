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

func (g *OutRoomsRepository) FindOneAndDeleteRoomUser(ctx context.Context, uID uint, RoomsID uint) error {
	result := g.GormDB.WithContext(ctx).Where("user_id = ? and Rooms_id = ?", uID, RoomsID).Delete(&mysql.RoomUsers{})
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), _errors.ErrRoomUserNotFound.Error(), utils.ErrFromClient)
	}
	return nil
}

func (g *OutRoomsRepository) FindOneAndUpdateRoom(ctx context.Context, RoomID uint) error {
	// 방 인원 -1
	result := g.GormDB.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count - 1"))
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)

	}
	return nil
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
