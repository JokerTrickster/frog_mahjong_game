package usecase

import (
	"context"
	"fmt"
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
	fmt.Println(ctx)

	// room create
	roomDTO := CreateRoomDTO(req, email)
	roomID, err := d.Repository.InsertOneRoom(ctx, roomDTO)
	if err != nil {
		return err
	}

	// room user create
	roomUserDTO := CreateRoomUserDTO(uID, roomID)
	err = d.Repository.InsertOneRoomUser(ctx, roomUserDTO)
	if err != nil {
		return err
	}

	return nil
}
