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
func (g *CreateRoomsRepository) FindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "play",
	}
	result := tx.WithContext(ctx).Model(user).Where("id = ?", uID).Updates(user)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), utils.HandleError(result.Error.Error(),uID,RoomID), utils.ErrFromClient)
	}
	return nil
}
func (g *CreateRoomsRepository) InsertOneRoom(ctx context.Context, tx *gorm.DB, RoomDTO mysql.Rooms) (int, error) {
	//방 인원이 최대 인원이 최소 인원보다 많거나 같고, 최대 인원이 2명 이상이거나 최소 인원이 2명 이상이어야 한다.
	if ((RoomDTO.MaxCount >= RoomDTO.MinCount) && (RoomDTO.MaxCount >= 2 || RoomDTO.MinCount >= 2)) == false {
		return 0, utils.ErrorMsg(ctx, utils.ErrUserNotFound, utils.Trace(), utils.HandleError(_errors.ErrBadRequest.Error(),RoomDTO), utils.ErrFromClient)
	}
	result := tx.WithContext(ctx).Create(&RoomDTO)
	if result.RowsAffected == 0 {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError("failed room insert one",RoomDTO), utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(result.Error.Error(),RoomDTO), utils.ErrFromMysqlDB)
	}
	return int(RoomDTO.ID), nil
}
func (g *CreateRoomsRepository) InsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomUserDTO mysql.RoomUsers) error {
	result := tx.WithContext(ctx).Create(&RoomUserDTO)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError("failed rooms user insert one",RoomUserDTO), utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(result.Error.Error(),RoomUserDTO), utils.ErrFromMysqlDB)
	}
	return nil
}
