package usecase

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"main/features/rooms/model/response"
	"main/utils/db/mysql"
	"time"

	"gorm.io/gorm"
)

type V02CreateRoomsUseCase struct {
	Repository     _interface.IV02CreateRoomsRepository
	ContextTimeout time.Duration
}

func NewV02CreateRoomsUseCase(repo _interface.IV02CreateRoomsRepository, timeout time.Duration) _interface.IV02CreateRoomsUseCase {
	return &V02CreateRoomsUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *V02CreateRoomsUseCase) V02Create(c context.Context, uID uint, email string, req *request.ReqV02Create) (response.ResV02CreateRoom, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	var res response.ResV02CreateRoom
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// Rooms V02Create
		RoomsDTO, err := V02CreateRoomDTO(ctx, req, uID)
		if err != nil {
			return err
		}
		RoomID, err := d.Repository.InsertOneRoom(ctx, tx, RoomsDTO)
		if err != nil {
			return err
		}

		// Rooms user V02Create
		RoomsUserDTO, err := V02CreateRoomUserDTO(uID, RoomID, "ready")
		if err != nil {
			return err
		}
		err = d.Repository.InsertOneRoomUser(ctx, tx, RoomsUserDTO)
		if err != nil {
			return err
		}

		// user 정보 변경 Rooms id와 state 변경
		err = d.Repository.FindOneAndUpdateUser(ctx, tx, uID, uint(RoomID))
		if err != nil {
			return err
		}

		res = response.ResV02CreateRoom{
			RoomID: RoomID,
		}
		return nil
	})

	return res, err
}
