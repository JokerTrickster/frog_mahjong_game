package usecase

import (
	"context"
	"fmt"
	_interface "main/features/room/model/interface"
	"main/features/room/model/request"
	"time"
)

type OutRoomUseCase struct {
	Repository     _interface.IOutRoomRepository
	ContextTimeout time.Duration
}

func NewOutRoomUseCase(repo _interface.IOutRoomRepository, timeout time.Duration) _interface.IOutRoomUseCase {
	return &OutRoomUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *OutRoomUseCase) Out(c context.Context, uID uint, req *request.ReqOut) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	fmt.Println(ctx)

	// roomID에 해당하는 userID를 삭제한다.
	err := d.Repository.FindOneAndDeleteRoomUser(ctx, uID, req.RoomID)
	if err != nil {
		return err
	}
	// room 현재 인원수를 -1한다.
	err = d.Repository.FindOneAndUpdateRoom(ctx, req.RoomID)
	if err != nil {
		return err
	}
	// user에 room_id를 1로 바꾸고 state를 wait으로 변경한다.
	err = d.Repository.FindOneAndUpdateUser(ctx, uID)
	if err != nil {
		return err
	}

	return nil
}
