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

type CreateRoomsUseCase struct {
	Repository     _interface.ICreateRoomsRepository
	ContextTimeout time.Duration
}

func NewCreateRoomsUseCase(repo _interface.ICreateRoomsRepository, timeout time.Duration) _interface.ICreateRoomsUseCase {
	return &CreateRoomsUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *CreateRoomsUseCase) Create(c context.Context, uID uint, email string, req *request.ReqCreate) (response.ResCreateRoom, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	var res response.ResCreateRoom
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// Rooms create
		RoomsDTO, err := CreateRoomDTO(ctx, req, uID)
		if err != nil {
			return err
		}
		RoomID, err := d.Repository.InsertOneRoom(ctx, tx, RoomsDTO)
		if err != nil {
			return err
		}

		res = response.ResCreateRoom{
			RoomID: RoomID,
		}
		return nil
	})

	return res, err
}
