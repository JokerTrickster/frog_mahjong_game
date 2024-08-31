package repository

import (
	"context"
	_errors "main/features/rooms/model/errors"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewJoinRoomsRepository(gormDB *gorm.DB) _interface.IJoinRoomsRepository {
	return &JoinRoomsRepository{GormDB: gormDB}
}

func (g *JoinRoomsRepository) FindOneRoom(ctx context.Context, tx *gorm.DB, req *request.ReqJoin) (mysql.Rooms, error) {
	// 방 참여 가능한지 체크
	RoomDTO := mysql.Rooms{}
	result := tx.WithContext(ctx).Where("id = ? and password = ?", req.RoomID, req.Password).First(&RoomDTO)
	if result.Error != nil {
		return mysql.Rooms{}, utils.ErrorMsg(ctx, utils.ErrRoomNotFound, utils.Trace(), _errors.ErrRoomNotFound.Error(), utils.ErrFromClient)
	}
	return RoomDTO, nil
}

func (g *JoinRoomsRepository) FindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) error {
	result := tx.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count + 1"))
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}

func (g *JoinRoomsRepository) FindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "play",
	}
	result := tx.WithContext(ctx).Model(&user).Where("id = ? and state = ?", uID, "wait").Updates(user)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}

	return nil
}

func (g *JoinRoomsRepository) InsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomUserDTO mysql.RoomUsers) error {
	result := tx.WithContext(ctx).Create(&RoomUserDTO)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), "failed room user insert one", utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
