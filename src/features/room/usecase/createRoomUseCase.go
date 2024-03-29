package usecase

import (
	"context"
	_interface "main/features/room/model/interface"
	"main/features/room/model/request"
	"time"
)

type CreateRoomUseCase struct {
	Repository     _interface.ICreateRoomRepository
	ContextTimeout time.Duration
}

func NewCreateRoomUseCase(repo _interface.ICreateRoomRepository, timeout time.Duration) _interface.ICreateRoomUseCase {
	return &CreateRoomUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *CreateRoomUseCase) Create(c context.Context, uID uint, email string, req *request.ReqCreate) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// room create
	roomDTO, err := CreateRoomDTO(ctx, req, email)
	if err != nil {
		return err
	}
	roomID, err := d.Repository.InsertOneRoom(ctx, roomDTO)
	if err != nil {
		return err
	}

	// room user create
	roomUserDTO, err := CreateRoomUserDTO(uID, roomID, "ready")
	if err != nil {
		return err
	}
	err = d.Repository.InsertOneRoomUser(ctx, roomUserDTO)
	if err != nil {
		return err
	}

	// user 정보 변경 room id와 state 변경
	err = d.Repository.FindOneAndUpdateUser(ctx, uID, uint(roomID))
	if err != nil {
		return err
	}

	return nil
}
