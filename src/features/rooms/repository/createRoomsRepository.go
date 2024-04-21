package repository

import (
	"context"
	_errors "main/features/rooms/model/errors"
	_interface "main/features/rooms/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewCreateRoomsRepository(gormDB *gorm.DB) _interface.ICreateRoomsRepository {
	return &CreateRoomsRepository{GormDB: gormDB}
}
func (g *CreateRoomsRepository) FindOneAndUpdateUser(ctx context.Context, uID uint, RoomID uint) error {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "play",
	}
	result := g.GormDB.WithContext(ctx).Model(user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), result.Error.Error(), utils.ErrFromClient)
	}
	return nil
}
func (g *CreateRoomsRepository) InsertOneRoom(ctx context.Context, RoomDTO mysql.Rooms) (int, error) {
	//방 인원이 최대 인원이 최소 인원보다 많거나 같고, 최대 인원이 2명 이상이거나 최소 인원이 2명 이상이어야 한다.
	if ((RoomDTO.MaxCount >= RoomDTO.MinCount) && (RoomDTO.MaxCount >= 2 || RoomDTO.MinCount >= 2)) == false {
		return 0, utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), _errors.ErrBadRequest.Error(), utils.ErrFromClient)
	}
	result := g.GormDB.WithContext(ctx).Create(&RoomDTO)
	if result.RowsAffected == 0 {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), "failed Room insert one", utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}
	return int(RoomDTO.ID), nil
}
func (g *CreateRoomsRepository) InsertOneRoomUser(ctx context.Context, RoomUserDTO mysql.RoomUsers) error {
	result := g.GormDB.WithContext(ctx).Create(&RoomUserDTO)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), "failed Rooms user insert one", utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
