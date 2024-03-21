package repository

import (
	"context"
	"fmt"
	_errors "main/features/room/model/errors"
	_interface "main/features/room/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewCreateRoomRepository(gormDB *gorm.DB) _interface.ICreateRoomRepository {
	return &CreateRoomRepository{GormDB: gormDB}
}
func (g *CreateRoomRepository) InsertOneRoom(ctx context.Context, roomDTO mysql.Rooms) (int, error) {
	//방 인원이 최대 인원이 최소 인원보다 많거나 같고, 최대 인원이 2명 이상이거나 최소 인원이 2명 이상이어야 한다.
	if ((roomDTO.MaxCount >= roomDTO.MinCount) && (roomDTO.MaxCount >= 2 || roomDTO.MinCount >= 2)) == false {
		return 0, utils.ErrorMsg(ctx, utils.ErrUserNotExist, utils.Trace(), _errors.ErrBadRequest.Error(), utils.ErrFromClient)
	}
	result := g.GormDB.WithContext(ctx).Create(&roomDTO)
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("failed room insert one")
	}
	if result.Error != nil {
		return 0, result.Error
	}
	return int(roomDTO.ID), nil
}
func (g *CreateRoomRepository) InsertOneRoomUser(ctx context.Context, roomUserDTO mysql.RoomUsers) error {
	result := g.GormDB.WithContext(ctx).Create(&roomUserDTO)
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed room user insert one")
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
