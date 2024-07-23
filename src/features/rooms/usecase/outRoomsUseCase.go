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
	roomDTO, err := d.Repository.FindOneAndUpdateRoom(ctx, req.RoomID)
	if err != nil {
		return err
	}
	// user에 rooms_id를 1로 바꾸고 state를 wait으로 변경한다.
	err = d.Repository.FindOneAndUpdateUser(ctx, uID)
	if err != nil {
		return err
	}

	if roomDTO.CurrentCount == 0 {
		// 방 삭제
		err = d.Repository.FindOneAndDeleteRoom(ctx, req.RoomID)
		if err != nil {
			return err
		}

	} else if roomDTO.CurrentCount == 1 {
		// 인원이 1명이면 남아 있는 유저를 방장으로 변경
		//방장이 나가면 다른 유저 중 한명을 방장으로 변경
		//룸에 남아있는 유저 정보를 가져온다.
		roomUserDTO, err := d.Repository.FindOneRoomUser(ctx, req.RoomID)
		if err != nil {
			return err
		}
		userDTO, err := d.Repository.FindOneUser(ctx, uint(roomUserDTO.UserID))
		if err != nil {
			return err
		}

		// 해당 유저를 방장으로 업데이트 한다.

		//방장으로 변경하기 위해 업데이트해야 될 부분들
		// rooms -> owner 변경
		err = d.Repository.ChangeRoomOnwer(ctx, req.RoomID, userDTO.ID)
		if err != nil {
			return err
		}

	}

	return nil
}
