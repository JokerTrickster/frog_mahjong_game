package usecase

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"main/features/rooms/model/response"
	"time"
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

	// Rooms create
	RoomsDTO, err := CreateRoomDTO(ctx, req, email)
	if err != nil {
		return response.ResCreateRoom{}, err
	}
	RoomID, err := d.Repository.InsertOneRoom(ctx, RoomsDTO)
	if err != nil {
		return response.ResCreateRoom{}, err
	}

	// Rooms user create
	RoomsUserDTO, err := CreateRoomUserDTO(uID, RoomID, "ready")
	if err != nil {
		return response.ResCreateRoom{}, err
	}
	err = d.Repository.InsertOneRoomUser(ctx, RoomsUserDTO)
	if err != nil {
		return response.ResCreateRoom{}, err
	}

	// user 정보 변경 Rooms id와 state 변경
	err = d.Repository.FindOneAndUpdateUser(ctx, uID, uint(RoomID))
	if err != nil {
		return response.ResCreateRoom{}, err
	}

	res := response.ResCreateRoom{
		RoomID: RoomID,
	}

	return res, nil
}
