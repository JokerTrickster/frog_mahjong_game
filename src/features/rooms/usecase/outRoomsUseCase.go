package usecase

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"time"
)

type OutRoomsUseCase struct {
	Repository     _interface.IOutRoomsRepository
	ContextTimeout time.Duration
}

func NewOutRoomsUseCase(repo _interface.IOutRoomsRepository, timeout time.Duration) _interface.IOutRoomsUseCase {
	return &OutRoomsUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *OutRoomsUseCase) Out(c context.Context, uID uint, req *request.ReqOut) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// RoomsID에 해당하는 userID를 삭제한다.
	err := d.Repository.FindOneAndDeleteRoomUser(ctx, uID, req.RoomID)
	if err != nil {
		return err
	}
	// Rooms 현재 인원수를 -1한다.
	err = d.Repository.FindOneAndUpdateRoom(ctx, req.RoomID)
	if err != nil {
		return err
	}
	// user에 Rooms_id를 1로 바꾸고 state를 wait으로 변경한다.
	err = d.Repository.FindOneAndUpdateUser(ctx, uID)
	if err != nil {
		return err
	}

	return nil
}
