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

func (g *JoinRoomsRepository) FindOneRoom(ctx context.Context, req *request.ReqJoin) (mysql.Rooms, error) {
	// 방 참여 가능한지 체크
	RoomDTO := mysql.Rooms{}
	result := g.GormDB.WithContext(ctx).Where("id = ? and password = ?", req.RoomID, req.Password).First(&RoomDTO)
	if result.Error != nil {
		return mysql.Rooms{}, utils.ErrorMsg(ctx, utils.ErrRoomNotFound, utils.Trace(), _errors.ErrRoomNotFound.Error(), utils.ErrFromClient)
	}
	return RoomDTO, nil
}

func (g *JoinRoomsRepository) FindOneAndUpdateRoom(ctx context.Context, RoomID uint) error {
	result := g.GormDB.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", RoomID).Update("current_count", gorm.Expr("current_count + 1"))
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}

func (g *JoinRoomsRepository) FindOneAndUpdateUser(ctx context.Context, uID uint, RoomID uint) error {
	user := mysql.Users{
		RoomID: int(RoomID),
		State:  "play",
	}
	result := g.GormDB.WithContext(ctx).Model(&user).Where("id = ? and state = ?", uID, "wait").Updates(user)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}

	return nil
}

func (g *JoinRoomsRepository) InsertOneRoomUser(ctx context.Context, RoomUserDTO mysql.RoomUsers) error {
	result := g.GormDB.WithContext(ctx).Create(&RoomUserDTO)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), "failed room user insert one", utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
