package usecase

import (
	"context"
	_errors "main/features/room/model/errors"
	_interface "main/features/room/model/interface"
	"main/features/room/model/request"
	"main/utils"
	"time"
)

type JoinRoomUseCase struct {
	Repository     _interface.IJoinRoomRepository
	ContextTimeout time.Duration
}

func NewJoinRoomUseCase(repo _interface.IJoinRoomRepository, timeout time.Duration) _interface.IJoinRoomUseCase {
	return &JoinRoomUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *JoinRoomUseCase) Join(c context.Context, uID uint, email string, req *request.ReqJoin) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 방 참여 가능한지 체크
	roomDTO, err := d.Repository.FindOneRoom(ctx, req)
	if err != nil {
		return err
	}
	if roomDTO.CurrentCount == roomDTO.MaxCount {
		return utils.ErrorMsg(ctx, utils.ErrRoomFull, utils.Trace(), _errors.ErrRoomFull.Error(), utils.ErrFromClient)
	}
	// 방 유저 정보를 생성한다.
	roomUserDTO, err := CreateRoomUserDTO(uID, int(req.RoomID), "wait")
	if err != nil {
		return err
	}
	err = d.Repository.InsertOneRoomUser(ctx, roomUserDTO)
	if err != nil {
		return err
	}
	// 방 현재 인원을 증가시킨다.
	err = d.Repository.FindOneAndUpdateRoom(ctx, req.RoomID)
	if err != nil {
		return err
	}

	//유저 정보를 업데이트 한다.
	err = d.Repository.FindOneAndUpdateUser(ctx, uID, req.RoomID)
	if err != nil {
		return err
	}

	return nil
}
