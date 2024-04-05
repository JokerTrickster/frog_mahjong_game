package repository

import (
	"context"
	_errors "main/features/room/model/errors"
	_interface "main/features/room/model/interface"
	"main/features/room/model/request"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewJoinRoomRepository(gormDB *gorm.DB) _interface.IJoinRoomRepository {
	return &JoinRoomRepository{GormDB: gormDB}
}

func (g *JoinRoomRepository) FindOneRoom(ctx context.Context, req *request.ReqJoin) (mysql.Rooms, error) {
	// 방 참여 가능한지 체크
	roomDTO := mysql.Rooms{}
	result := g.GormDB.WithContext(ctx).Where("id = ? and password = ?", req.RoomID, req.Password).First(&roomDTO)
	if result.Error != nil {
		return mysql.Rooms{}, utils.ErrorMsg(ctx, utils.ErrRoomImpossibleJoin, utils.Trace(), _errors.ErrRoomNotFound.Error(), utils.ErrFromClient)
	}
	return roomDTO, nil
}

func (g *JoinRoomRepository) FindOneAndUpdateRoom(ctx context.Context, roomID uint) error {
	result := g.GormDB.WithContext(ctx).Model(&mysql.Rooms{}).Where("id = ?", roomID).Update("current_count", gorm.Expr("current_count + 1"))
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}

func (g *JoinRoomRepository) FindOneAndUpdateUser(ctx context.Context, uID uint, roomID uint) error {
	user := mysql.Users{
		RoomID: int(roomID),
		State:  "play",
	}
	result := g.GormDB.WithContext(ctx).Model(&user).Where("id = ? and state = ?", uID, "wait").Updates(user)
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}

	return nil
}

func (g *JoinRoomRepository) InsertOneRoomUser(ctx context.Context, roomUserDTO mysql.RoomUsers) error {
	result := g.GormDB.WithContext(ctx).Create(&roomUserDTO)
	if result.RowsAffected == 0 {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), "failed room user insert one", utils.ErrFromMysqlDB)
	}
	if result.Error != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), result.Error.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
