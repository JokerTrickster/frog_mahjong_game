package repository

import (
	"context"
	_interface "main/features/room/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewOutRoomRepository(gormDB *gorm.DB) _interface.IOutRoomRepository {
	return &OutRoomRepository{GormDB: gormDB}
}

func (g *OutRoomRepository) FindOneAndDeleteRoomUser(ctx context.Context, uID uint, roomID uint) error {
	result := g.GormDB.WithContext(ctx).Where("user_id = ? and room_id = ?", uID, roomID).Delete(&mysql.RoomUsers{})
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)
	}
	return nil
}

func (g *OutRoomRepository) FindOneAndUpdateRoom(ctx context.Context, roomID uint) error {
	// 방 인원 -1
	result := g.GormDB.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", roomID).Update("current_count", gorm.Expr("current_count - 1"))
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)

	}
	return nil
}

func (g *OutRoomRepository) FindOneAndUpdateUser(ctx context.Context, uID uint) error {
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
